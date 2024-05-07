#![feature(write_all_vectored)]

use std::thread;
use std::time::Duration;
use std::io::{self, BufWriter, Write};
use console::{Term};

fn main() {

        let mut term = Term::stdout();
        let mut sw = BufWriter::new(io::stdout());
        let mut buffer = [[' ' as u8; 80]; 25];
        let ramp = ".:-=+*#%@".as_bytes();
        let mut t = 0.0f32;
        loop {
            let mut sy = 0;
            let s = t.sin();
            let a = s * s * s * s * 0.2f32;
            let mut z = 1.3f32;
            while z > -1.2f32 {
                let mut py: usize = 0;
                let tz = z * (1.2f32 - a);
                let mut x = -1.5f32;
                while x < 1.5f32 {
                    let tx = x * (1.2f32 + a);
                    let v = f(tx, 0.0f32, tz);
                    if v <= 0.0f32 {
                        let y0 = h(tx, tz);
                        let ny = 0.01f32;
                        let nx = h(tx + ny, tz) - y0;
                        let nz = h(tx, tz + ny) - y0;
                        let nd = 1.0f32 / (nx * nx + ny * ny + nz * nz).sqrt();
                        let d = (nx + ny - nz) * nd * 0.5f32 + 0.5f32;
                        buffer[sy][py] = ramp[(d * 5.0f32) as usize];
                    } else {
                        buffer[sy][py] = ' ' as u8;
                    }
                    x += 0.05f32;
                    py += 1;
                }
                z -= 0.1f32;
                sy+=1;
            }

            for sy in 0..25 {
                term.move_cursor_to(0,sy as usize);
                term.write(&buffer[sy as usize]);
            }


            thread::sleep(Duration::from_millis(16));
            t += 0.1f32;
        }
}

// IDEA bug x * x 这里会报错Use of moved value
fn f(x: f32, y: f32, z: f32) -> f32 {
    let a = x * x + 9.0f32 / 4.0f32 * y * y + z * z - 1.0;
    return a * a * a - x * x * z * z * z - 9.0f32 / 80.0f32 * y * y * z * z * z;
}

fn h(x: f32, z: f32) -> f32 {
    let mut y = 1.0f32;
    while y >= 0.0f32 {
        y -= 0.001f32;
        if f(x, y, z) <= 0.0f32 {
            return y;
        }
    }
    return 0.0f32;
}