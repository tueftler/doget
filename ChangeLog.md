DoGet change log
================

## ?.?.? / ????-??-??

## 0.8.0 / 2016-09-08

* Added `clean` subcommand as requested in #9. See pull request #13
  (@mikey179, @thekid)
* Fixed issue #15: Escaped character leads to premature end of line.
  See pull request #16
  (@mikey179)

## 0.7.0 / 2016-08-30

* Fixed issue with paths in `ADD` and `COPY` instructions not being
  correctly resolved to vendor directory
  (@thekid, @kiesel)

## 0.6.0 / 2016-08-29

* Fixed nil pointer dereference when handling nonexistant repositories
  (@thekid)
* Fixed leading and trailing whitespace in trait references leading to
  downloads breaking
  (@thekid, @mikey179)

## 0.5.0 / 2016-08-26

* Added `Emit()` function to all dockerfile Statement instances
  (@thekid)

## 0.4.0 / 2016-08-25

* Merged PR #6: Make default configuration builtin instead of shipping it as 
  a file. Implementes feature request #5
  (@thekid)
* Changed `config.Merge()` to only parse given files once - @mikey179 

## 0.3.0 / 2016-08-23

* Fixed issue #2: Added support for bitbucket.org downloads
  (@thekid)
* Fixed issue #4: Panic "invalid memory address or nil pointer dereference"
  (@thekid)

## 0.2.0 / 2016-08-22

* Added support for user config files in HOME (Un*x) and APPDATA (Windows)
  See issue #3 (*still missing XDG compliance!*)
  (@thekid)

## 0.1.0 / 2016-08-22

* Hello World! First release - @thekid