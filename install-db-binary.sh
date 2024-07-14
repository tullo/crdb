#!/bin/bash

BINARY=cockroach-v24.1.2.linux-amd64
TEMPDIR=/tmp/crdb

echo "Installing $BINARY"
mkdir -p $TEMPDIR /usr/local/lib/cockroach &&
wget -O $TEMPDIR/crdb.tgz https://binaries.cockroachdb.com/$BINARY.tgz &&
tar -zxf $TEMPDIR/crdb.tgz -C $TEMPDIR &&
install $TEMPDIR/$BINARY/cockroach /usr/local/bin/ &&
install $TEMPDIR/$BINARY/lib/libgeos.so /usr/local/lib/cockroach &&
install $TEMPDIR/$BINARY/lib/libgeos_c.so /usr/local/lib/cockroach &&
rm -rf $TEMPDIR

cockroach version
