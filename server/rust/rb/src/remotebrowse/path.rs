//! 远端路径解析：Windows 仅接受 `D:/…` / `\\…` 形态；Unix 接受 `/…`。
use std::borrow::Cow;
use std::path::{Component, Path, PathBuf};

pub fn default_browse_root() -> Result<PathBuf, String> {
    #[cfg(windows)]
    {
        windows_user_profile_dir().ok_or_else(|| {
            "无法解析 Windows 用户目录（请设置 USERPROFILE、PROFILE 或 HOMEDRIVE+HOMEPATH）".into()
        })
    }
    #[cfg(not(windows))]
    {
        if let Ok(v) = std::env::var("HOME") {
            let v = v.trim();
            if !v.is_empty() {
                let p = normalize_client_path(v)?;
                if p.is_dir() {
                    return Ok(p);
                }
            }
        }
        std::env::current_dir().map_err(|e| e.to_string())
    }
}

/// Windows 用户目录：仅读 `USERPROFILE` / `PROFILE` / `HOMEDRIVE`+`HOMEPATH` 环境变量。
#[cfg(windows)]
fn windows_user_profile_dir() -> Option<PathBuf> {
    for key in ["USERPROFILE", "PROFILE"] {
        if let Ok(v) = std::env::var(key) {
            if let Some(p) = normalize_windows_profile_env(&v) {
                if p.is_dir() {
                    return Some(p);
                }
            }
        }
    }
    if let (Ok(drive), Ok(home)) = (std::env::var("HOMEDRIVE"), std::env::var("HOMEPATH")) {
        let combined = format!("{}{}", drive.trim(), home.trim());
        if let Some(p) = normalize_windows_profile_env(&combined) {
            if p.is_dir() {
                return Some(p);
            }
        }
    }
    None
}

/// MSYS 可能把 USERPROFILE 设为 `/c/Users/…`，仅用于读取环境变量时转为 `C:/Users/…`。
#[cfg(windows)]
fn normalize_windows_profile_env(raw: &str) -> Option<PathBuf> {
    let t = raw.trim();
    if t.is_empty() {
        return None;
    }
    if let Some(w) = msys_drive_slash_to_win_str(t) {
        return normalize_client_path(&w).ok().filter(|p| !p.as_os_str().is_empty());
    }
    normalize_client_path(t).ok()
}

#[cfg(windows)]
fn msys_drive_slash_to_win_str(s: &str) -> Option<String> {
    let s = s.replace('\\', "/");
    if s.len() >= 4 && s.starts_with('/') {
        let b = s.as_bytes();
        if b[1].is_ascii_alphabetic() && b[2] == b'/' {
            let drive = (b[1] as char).to_ascii_uppercase();
            return Some(format!("{drive}:{}", &s[2..]));
        }
    }
    None
}

/// ReadFile 请求 path 后缀：客户端拉原图时追加，解析前剥离。
pub const READ_ORIGINAL_SUFFIX: &str = "#rb-original";

pub fn strip_read_original_suffix(client: &str) -> (Cow<'_, str>, bool) {
    if let Some(base) = client.strip_suffix(READ_ORIGINAL_SUFFIX) {
        return (Cow::Borrowed(base), true);
    }
    (Cow::Borrowed(client), false)
}

pub fn resolve_viewer_path(sandbox: Option<&str>, client: &str) -> Result<PathBuf, String> {
    let client = strip_windows_verbatim_prefix(client.trim());
    if !client.is_empty() && client != "." {
        ensure_client_path_format(&client)?;
    }
    let base = browse_base(sandbox)?;
    let mut path = if client.is_empty() || client == "." {
        base.clone()
    } else {
        normalize_client_path(client.as_ref())?
    };
    if path.is_relative() {
        path = base.join(path);
    }
    let path = normalize_path(&path);
    if let Some(sb) = sandbox.filter(|s| !s.is_empty()) {
        enforce_sandbox(&path, sb)?;
    }
    Ok(path)
}

pub fn display_path(path: &Path) -> String {
    to_wire_path(&canonical_or_identity(path))
}

/// 对外路径统一为 `D:/…`（Windows `\` → `/`）。
pub fn to_wire_path(path: &Path) -> String {
    let s = strip_windows_verbatim_prefix(&path.display().to_string()).into_owned();
    to_wire_path_str(&s)
}

fn to_wire_path_str(s: &str) -> String {
    match normalize_client_path(s) {
        Ok(p) => p.display().to_string().replace('\\', "/"),
        Err(_) => s.replace('\\', "/"),
    }
}

pub fn normalize_client_path(s: &str) -> Result<PathBuf, String> {
    let raw = strip_windows_verbatim_prefix(s.trim());
    let s = raw.replace('\\', "/");
    ensure_client_path_format(&s)?;
    let s = s.trim_end_matches('/');
    if s.is_empty() {
        return Ok(PathBuf::new());
    }
    if s.len() == 2 {
        let b = s.as_bytes();
        if b[0].is_ascii_alphabetic() && b[1] == b':' {
            return Ok(PathBuf::from(format!("{s}/")));
        }
    }
    Ok(PathBuf::from(s))
}

fn browse_base(sandbox: Option<&str>) -> Result<PathBuf, String> {
    if let Some(sb) = sandbox.filter(|s| !s.is_empty()) {
        return normalize_client_path(sb);
    }
    default_browse_root()
}

/// Windows wire 路径须为 `D:/…` 或 `//server/share`；Unix 拒绝 `C:/…` 盘符路径。
fn ensure_client_path_format(s: &str) -> Result<(), String> {
    let t = s.trim();
    if t.is_empty() || t == "." {
        return Ok(());
    }
    let norm = t.replace('\\', "/");
    #[cfg(windows)]
    {
        if norm.starts_with("//") {
            return Ok(());
        }
        if norm.starts_with('/') {
            return Err("Windows 路径请使用 D:/… 或 \\\\server\\share 格式".into());
        }
        return Ok(());
    }
    #[cfg(not(windows))]
    {
        if norm.len() >= 2 {
            let b = norm.as_bytes();
            if b[0].is_ascii_alphabetic() && b[1] == b':' {
                return Err("请使用 Unix 绝对路径，如 /home/user/…".into());
            }
        }
        Ok(())
    }
}

fn normalize_path(path: &Path) -> PathBuf {
    let mut out = PathBuf::new();
    for c in path.components() {
        match c {
            Component::ParentDir => {
                out.pop();
            }
            Component::CurDir => {}
            other => out.push(other),
        }
    }
    out
}

fn enforce_sandbox(path: &Path, sandbox: &str) -> Result<(), String> {
    let sb = normalize_client_path(sandbox)?;
    let sb = canonical_or_identity(&sb);
    let p = canonical_or_identity(path);
    if path_under_sandbox(&p, &sb) {
        return Ok(());
    }
    Err(format!(
        "path outside RB_AGENT_SANDBOX (path={}, sandbox={})",
        to_wire_path(&p),
        to_wire_path(&sb)
    ))
}

fn canonical_or_identity(path: &Path) -> PathBuf {
    path.canonicalize().unwrap_or_else(|_| normalize_path(path))
}

fn path_under_sandbox(path: &Path, sandbox: &Path) -> bool {
    #[cfg(windows)]
    {
        let p = path_key(path);
        let s = path_key(sandbox);
        return p == s || p.starts_with(&format!("{s}/"));
    }
    #[cfg(not(windows))]
    {
        path == sandbox || path.starts_with(sandbox)
    }
}

#[cfg(windows)]
fn path_key(path: &Path) -> String {
    path.display().to_string().replace('\\', "/").to_ascii_lowercase()
}

pub fn strip_windows_verbatim_prefix(s: &str) -> Cow<'_, str> {
    if let Some(rest) = s.strip_prefix(r"\\?\") {
        if let Some(unc) = rest.strip_prefix("UNC\\") {
            return Cow::Owned(format!(r"\\{unc}"));
        }
        return Cow::Borrowed(rest);
    }
    if let Some(rest) = s.strip_prefix("//?/") {
        return Cow::Borrowed(rest);
    }
    if let Some(rest) = s.strip_prefix("/?/") {
        return Cow::Borrowed(rest);
    }
    Cow::Borrowed(s)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn win_drive_only() {
        assert_eq!(normalize_client_path("D:").unwrap(), PathBuf::from("D:/"));
    }

    #[test]
    fn win_path_slash() {
        assert_eq!(
            normalize_client_path(r"D:\Users\lbyi").unwrap(),
            PathBuf::from("D:/Users/lbyi")
        );
        assert_eq!(to_wire_path_str(r"D:\code\hopeio"), "D:/code/hopeio");
    }

    #[cfg(windows)]
    #[test]
    fn msys_userprofile_env_to_win() {
        assert_eq!(
            normalize_windows_profile_env("/c/Users/lbyi").unwrap(),
            PathBuf::from("C:/Users/lbyi")
        );
        assert_eq!(
            normalize_windows_profile_env(r"C:\Users\lbyi").unwrap(),
            PathBuf::from("C:/Users/lbyi")
        );
    }

    #[cfg(windows)]
    #[test]
    fn reject_unix_style_input() {
        assert!(normalize_client_path("/d/code/hopeio").is_err());
        assert!(normalize_client_path("/home/lbyi").is_err());
        assert!(resolve_viewer_path(None, "/d/code").is_err());
    }

    #[test]
    fn sandbox_resolves_win_client_path() {
        let sb = "D:/Users/lbyi";
        let p = resolve_viewer_path(Some(sb), "D:/Users/lbyi/Documents").unwrap();
        assert_eq!(to_wire_path(&p), "D:/Users/lbyi/Documents");
    }

    #[test]
    fn read_original_suffix_stripped() {
        let (p, orig) = strip_read_original_suffix("D:/a.jpg#rb-original");
        assert!(orig);
        assert_eq!(p, "D:/a.jpg");
        let (p2, orig2) = strip_read_original_suffix("D:/a.jpg");
        assert!(!orig2);
        assert_eq!(p2, "D:/a.jpg");
    }
}
