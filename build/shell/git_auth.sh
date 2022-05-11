git filter-branch -f --env-filter '
if [ "$GIT_COMMITTER_NAME" = "oldname" ];
then
    GIT_COMMITTER_NAME="newname";
    GIT_COMMITTER_EMAIL="newaddr";
    GIT_AUTHOR_NAME="newname";
    GIT_AUTHOR_EMAIL="newaddr";
fi

if [ "$GIT_AUTHOR_NAME" = "oldname" ];
then
    GIT_COMMITTER_NAME="newname";
    GIT_COMMITTER_EMAIL="newaddr";
    GIT_AUTHOR_NAME="newname";
    GIT_AUTHOR_EMAIL="newaddr";
fi
' -- --all


git config --global http.proxy 'socks5://127.0.0.1:1080'

git config --global https.proxy 'socks5://127.0.0.1:1080'

vi ~/.gitconfig
[http]
proxy = socks5://127.0.0.1:2080
[https]
proxy = socks5://127.0.0.1:2080