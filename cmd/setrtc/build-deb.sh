#!/bin/bash

RELEASE_DATE=$(date -u +'%Y%m%d')

OPTIONS+=" --type=debian"
OPTIONS+=" --default"     # answer yes to any questions
OPTIONS+=" --fstrans"     # simulate filesystem so root permissions not needed
OPTIONS+=" --install=no"  # just make, don't install the package
OPTIONS+=" --nodoc"       # as long as actual docs don't exist

OPTIONS+=" --pkgname=mist-setrtc"
OPTIONS+=" --pkgversion=0.1.0"
OPTIONS+=" --pkgrelease=$RELEASE_DATE"
OPTIONS+=" --provides=mist-setrtc"
OPTIONS+=" --pkglicense=MIT"
OPTIONS+=" --maintainer=somebody@thinnect.com"

checkinstall $OPTIONS
