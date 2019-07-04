#/bin/bash
#
# This contains a few things I wanted to actually see running

PREFIX=feature/
git_local_branches() { git branch --no-color | sed 's/^[* ] //'; }
feature_branches=$(echo "$(git_local_branches)" | grep "^$PREFIX")
short_names=$(echo "$feature_branches" | sed "s ^$PREFIX  g")

max() { [[ "$1" -ge "$2" ]] && echo "$1" || echo "$2"; }

for branch in $short_names; do
    len=${#branch}
    echo $len
    width=$(max $width $len)
    echo $width
done
width=$(($width+5))
echo $width
for branch in $short_names; do
    printf "%-${width}s" "$branch"
done
