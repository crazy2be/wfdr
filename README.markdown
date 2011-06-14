WFDR Framework - Beta Release
=============================

Introduction
------------
There's a million different web frameworks out there, each for a different language. Ruby has rails, java has grails (and plenty of others), python has django. However, although each of these frameworks has their own merits, they are designed for a single language, and a single language only. As of yet, I have been unable to locate a web framework that works with go (the excellent web.go is a library, not a framework). Despite no framework existing, I had a website to write, and the builtin HTTP library and a bad experience with PHP lead me to choose go as my language of choice. The framework has evolved along with the website, and has been in development for almost a year now. However, the initial builds were not really a framework at all, and only in the past few months has it begun to look viable for others to use. Many things are still broken or could use work, and, only recently learning of gotest, unit tests are needed for much of the code. Most of this should get cleaned up over the summer, especially with your help. Think of this framework as a starting point rather than a final, polished product.

tl;dr: Framework designed for go, could work with other languages, not production ready.

Philosophy
----------
The WFDR framework is designed around the principle of least privilege, separating websites into individual "modules". Each module has a completely independent source folder in `modules/`, and is placed and run from its own jail in `jails/`. Eventually, modules will be chrooted into this jail in a production configuration. Each module represents a subset of the functionality available on the site, and is represented on the client-side by a different top-level "folder". For example, a website might have a news module, available at `/news/*`, and a photos module, available at `/photos/*`. Requests to the folder or subfolders thereof automatically get forwarded to the module.

The framework is also heavily inspired by git and unix designs, opting to use many different commands rather than one monolithic binary. Indeed, the framework itself arose because the design of a single binary for all server functionality was broken by design, providing a single point of failure that was very likely to fail. A memory leak in any server-side code would, in a monolithic design, cause the entire server to be affected or stop responding.

Getting Started
---------------
The WFDR framework is designed to be simple to get started with, and should work fairly well with whatever your favorite HTTP library is.

    git clone git://github.com/crazy2be/wfdr.git
    ./compile

If all goes well, you've now built the framework. If there is an error while compiling, please let me know by filing a bug report. The compile script should ensure everything it needs is installed.

Now, to start the framework, open up two terminals. In both terminals, `cd` to the location where you installed the framework. In the first terminal, run `wfdr-deamon` (sic). This manages your various module processes, as well as starting some other programs to make sure that files are synced automatically when changed (on Linux at least). In the second terminal, you can now use the "wfdr" command to control modules by communicating with the daemon. Let's start a few of the included example modules:

    wfdr start base main auth photos pages news

If all goes well, you should now be able to navigate to http://localhost:8080/ to see these modules running.

Writing Your Own Module(s)
--------------------------
The WFDR framework is designed to make writing additional modules as bits of functionality as easy as possible. Each module is located in `modules/`, and includes everything needed for the module to properly function, including css, js, images, (mustache) templates, and source code. You can look at some of the other modules to get an idea of how things work, but they can really work any way you wish, as long as you do the following:

 -  Put your source files in `modules/<name>/src`, and have an accompanying Makefile that will compile the source files.
 -  The compiled executable should go in `modules/<name>/bin`. This executable should have the same name as the module, making the full path `modules/<name>/bin/<name>`. Alternatively, you can name the executable anything you want (it should still go in `modules/<name>/bin`), and have a shell script named with the same name as the module in `modules/<name>/sh`. If this script has +x set, it will be run instead.