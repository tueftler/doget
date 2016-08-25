DoGet change log
================

## ?.?.? / ????-??-??

## 0.4.0 / 2014-08-25

* Merged PR #6: Make default configuration builtin instead of shipping it as 
  a file. Implementes feature request #5
  (@thekid)
* Changed `config.Merge()` to only parse given files once - @mikey179 

## 0.3.0 / 2014-08-23

* Fixed issue #2: Added support for bitbucket.org downloads
  (@thekid)
* Fixed issue #4: Panic "invalid memory address or nil pointer dereference"
  (@thekid)

## 0.2.0 / 2014-08-22

* Added support for user config files in HOME (Un*x) and APPDATA (Windows)
  See issue #3 (*still missing XDG compliance!*)
  (@thekid)

## 0.1.0 / 2014-08-22

* Hello World! First release - @thekid