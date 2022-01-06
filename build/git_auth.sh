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