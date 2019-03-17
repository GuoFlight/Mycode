# .bashrc

# User specific aliases and functions

alias rm='rm -i'
alias cp='cp -i'
alias mv='mv -i'

#This is my alias
alias sys='systemctl'
alias dcoker='docker'
alias 

PATH="$PATH:/root"

# Source global definitions
if [ -f /etc/bashrc ]; then
	. /etc/bashrc
fi
