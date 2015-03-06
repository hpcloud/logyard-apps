#
# Makefile for stackato-logyard-apps
#
# Used solely by packaging systems.
# Must support targets "all", "install", "uninstall".
#
# During the packaging install phase, the native packager will
# set either DESTDIR or prefix to the directory which serves as
# a root for collecting the package files.
#
# Additionally, stackato-pkg sets STACKATO_PKG_BRANCH to the
# current git branch of this package, so that we may use it to
# fetch other git repos with the corresponding branch.
#
# The resulting package installs in /home/stackato,
# is not intended to be relocatable.
#
# To locally test this Makefile, run:
#
#   rm -rf .gopath; STACKATO_PKG_BRANCH=mybranch make
#

NAME=github.com/ActiveState/logyard-apps

SRCDIR=src/$(NAME)

INSTALLHOME=/home/stackato
INSTALLROOT=$(INSTALLHOME)/stackato
GOBINDIR=$(INSTALLROOT)/go/bin

INSTDIR=$(DESTDIR)$(prefix)

INSTHOMEDIR=$(INSTDIR)$(INSTALLHOME)
INSTROOTDIR=$(INSTDIR)$(INSTALLROOT)
INSTGOPATH=$(INSTDIR)$(INSTALLROOT)/go
INSTBINDIR=$(INSTDIR)$(INSTALLHOME)/bin

BUILDGOPATH=$(shell pwd)/.gopath

GOARGS=-v -tags zmq_3_x

GOARGS_TEST=-race

export PATH := /usr/local/go/bin:$(BUILDGOPATH)/bin/:$(PATH)

all:	repos compile

repos:
	mkdir -p $(BUILDGOPATH)/src/$(NAME)
	git archive HEAD | tar -x -C $(BUILDGOPATH)/src/$(NAME)
	GOPATH=$(BUILDGOPATH) GOROOT=/usr/local/go go get -v github.com/vube/depman
	GOPATH=$(BUILDGOPATH) GOROOT=/usr/local/go depman
	rm -f $(BUILDGOPATH)/bin/depman

compile:	$(BUILDGOROOT)
	GOPATH=$(BUILDGOPATH) GOROOT=/usr/local/go go install $(GOARGS) $(NAME)/...

install:
	mkdir -p $(INSTGOPATH)/$(SRCDIR)
	rsync -a $(BUILDGOPATH)/$(SRCDIR)/etc/*.yml $(INSTGOPATH)/$(SRCDIR)/etc/
	rsync -a $(BUILDGOPATH)/bin $(INSTGOPATH)
	rsync -a etc $(INSTROOTDIR)
	mkdir -p $(INSTBINDIR)
	chown -Rh stackato.stackato $(INSTHOMEDIR)

clean:	$(BUILDGOROOT)
	GOPATH=$(BUILDGOPATH) GOROOT=/usr/local/go go clean

# For developer use.

all-local: fmt repos-local compile-local

repos-local:
	mkdir -p $(BUILDGOPATH)/src/$(NAME)
	git archive HEAD | tar -x -C $(BUILDGOPATH)/src/$(NAME)
	GOPATH=$(BUILDGOPATH) GOROOT=${GOROOT} go get -v github.com/vube/depman
	GOPATH=$(BUILDGOPATH) GOROOT=${GOROOT} depman
	rm -f $(BUILDGOPATH)/bin/depman


compile-local:	$(BUILDGOROOT)
	GOPATH=$(BUILDGOPATH) GOROOT=${GOROOT} go install $(GOARGS) $(NAME)/...

dev-install:	fmt dev-installall

# convenient alias
i:	dev-install

dev-installall:
	go install $(GOARGS) $(NAME)/...

fmt:
	gofmt -w .

dev-test:
	go test $(GOARGS) $(GOARGS_TEST) $(NAME)/...