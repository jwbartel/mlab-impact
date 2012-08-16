#!/usr/bin/make -f
#
# Copyright 2011 Google Inc. All Rights Reserved.
# Author: gavaletz@google.com (Eric Gavaletz)


# BUILD REQUIREMENTS

# java -- sudo apt-get install openjdk-7-jre
# TODO(user) add appcfg.py, dev_appserver.py and gjslint to your PATH.
# GAE -- http://code.google.com/appengine/
# gjslint -- http://code.google.com/p/closure-linter/

# http://code.google.com/p/closure-library/
closure_root = $(HOME)/SRC/closure-library/
builder = $(closure_root)closure/bin/build/closurebuilder.py

# http://code.google.com/p/closure-compiler/
compiler_jar = $(HOME)/SRC/closure-compiler/compiler.jar

# http://code.google.com/p/closure-stylesheets/
#TODO(user) you will need to create a symbolic link for closure-stylesheets.jar.
compiler_css = java -jar $(HOME)/SRC/closure-stylesheets/closure-stylesheets.jar


project_root= ./
GAE_STATIC = gae/mlab-impact/static/
version = $(shell egrep "^version: " gae/mlab-impact/app.yaml | cut -d " " -f 2)


cfb = --compiler_flags=
compiler_flags = $(cfb)"--compilation_level=ADVANCED_OPTIMIZATIONS"
compiler_flags += $(cfb)"--externs=externs/google_analytics_api.js"
compiler_flags += $(cfb)"--externs=externs/google_loader_api.js"
compiler_flags += $(cfb)"--externs=externs/google_maps_api_v3_6.js"
compiler_flags += $(cfb)"--externs=externs/google_visualization_api.js"
# Uncomment these for debuging.
#compiler_flags += $(cfb)"--formatting=PRETTY_PRINT"
#compiler_flags += $(cfb)"--debug"


compile = $(builder) \
	--root=$(closure_root) \
	--root=$(project_root) \
	--output_mode=compiled \
	--compiler_jar=$(compiler_jar) \
	$(compiler_flags) \
	--namespace=NAMESPACE \
	--output_file=


depends = $(builder) \
	--root=$(closure_root) \
	--root=$(project_root) \
	--output_mode=list \
	--compiler_jar=$(compiler_jar) \
	--namespace=NAMESPACE \


.PHONY: all
all: impact css images


.PHONY: impact
impact: $(GAE_STATIC)$(version).impact-comp.js


.PHONY: css
css: $(GAE_STATIC)$(version).impact-comp.css


$(GAE_STATIC)$(version).impact-comp.js: $(shell $(depends:NAMESPACE=impact))
	$(compile:NAMESPACE=impact)$(GAE_STATIC)$(version).impact-comp.js


$(GAE_STATIC)$(version).impact-comp.css: css/*.css
	$(compiler_css) --output-file $(GAE_STATIC)$(version).impact-comp.css css/*.css


$(GAE_STATIC)robots.txt: robots.txt
	cp robots.txt $(GAE_STATIC)


.PHONY: images
images:
	rsync -av images $(GAE_STATIC)
