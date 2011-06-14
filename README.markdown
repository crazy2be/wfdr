WFDR Framework - Beta Release
=============================

Introduction
------------
There's a million different web frameworks out there, each for a different language. Ruby has rails, java has grails (and plenty of others), python has django. However, although each of these frameworks has their own merits, they are designed for a single language, and a single language only. As of yet, I have been unable to locate a web framework that works with go (the excellent web.go is a library, not a framework). Despite no framework existing, I had a website to write, and the builtin HTTP library and a bad experience with PHP lead me to choose go as my language of choice. The framework has evolved along with the website, and has been in development for almost a year now. However, the initial builds were not really a framework at all, and only in the past few months has it begun to look viable for others to use. Many things are still broken or could use work, and, only recently learning of gotest, unit tests are needed for much of the code. Most of this should get cleaned up over the summer, especially with your help. Think of this framework as a starting point rather than a final, polished product.

tl;dr: Framework designed for go, could work with other languages, not production ready.

Philosophy
----------
The WFDR framework is designed around the principle of least privilege, separating websites into individual "modules". Each module has a completely independent source folder in `modules/`, and is placed and run from its own jail in `jails/`. Eventually, modules will be chrooted into this jail in a production configuration. Each module represents a subset of the functionality available on the site, and is represented on the client-side by a different top-level "folder". For example, a website might have a news module, available at /news/*, and a photos module, available at /photos/*. Requests to the folder or subfolders thereof automatically get forwarded to the module.

The framework is also heavily inspired by git and unix designs, opting to use many different commands rather than one monolithic binary. Indeed, the framework itself arose because the design of a single binary for all server functionality was broken by design, providing a single point of failure that was very likely to fail. A memory leak in any server-side code would, in a monolithic design, cause the entire server to be affected or stop responding.