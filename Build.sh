#!/bin/sh
watcher -cmd="sh Update.sh" -recursive -pipe=true -list ./control &
canthefason_watcher -run dogfooter-control/control/cmd -watch dogfooter-control
