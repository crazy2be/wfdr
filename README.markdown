WFDR Framework - Beta Release
=============================

Background
----------
There's a million different web frameworks out there, each for a different language. Ruby has rails, java has grails, python has django. However, although each of these frameworks has many merits, they are designed for a single language, and none of them work with go. Regardless, I had a website to write, and go's strengths in that area lead me to chose it as my language of choice. The framework has evolved along with the website, and has been in development for almost a year now. Although the initial builds were not really a framework at all, in the past few months it has evolved to become quite viable for others to use. As a result, I am publishing the source in the hope that others many find it useful, perhaps giving me some constructive feedback and criticism along the way. Expect a few rough edges.

**tl;dr:** Web framework designed for go, has the ability to work with other languages, not production ready.

Features
--------

 -  **Serves Files**: Static (or dynaically generated!) js, css, and image files are automatically served by the framework for you.
 -  **Customizable Layouts**: Customize as little or as much as you want between mobile and desktop clients, with more to come. [(More on layouts)](https://github.com/crazy2be/wfdr/wiki/Layouts)
 -  **Language Agnostic**: Although several features only have bindings for go, implementing them for other languages would be trivial.
 -  **Encourages Modular Design**: Each *Module* consists of an isolated and separated section of a site. This separation occurs both in the source code ([more on modules](https://github.com/crazy2be/wfdr/wiki/Modules)), and at runtime ([more on jails](https://github.com/crazy2be/wfdr/wiki/Jails)). A bug in one module will leave other modules completely unaffected.
 -  **Customizable**: Easy to understand and hack, thanks to a design inspired by UNIX and git.

Getting Started
---------------

    git clone git://github.com/crazy2be/wfdr.git
    cd wfdr
    ./compile

If all goes well, you've now built the framework. If there is an error while compiling, please let me know by filing a bug report. The compile script should ensure everything it needs is installed.

Now, to start the framework, open up two terminals. In both terminals, `cd` to the location where you installed the framework. In the first terminal, run `wfdr-deamon` (sic). This daemon process manages your various module processes, as well as starting some other programs to make sure that files are synced automatically when changed (uses inotify; Linux only atm). In the second terminal, you can now use the `wfdr` command to control modules by communicating with the daemon. Let's start a few of the included example modules:

    wfdr start base main auth photos pages news

If all goes well, you should now be able to navigate to http://localhost:8080/ and see these modules running. If not, that's a bug, either in the framework or the documentation. Let me know.

What's Next?
------------
The WFDR framework is designed to be simple to get started with, and should work with whatever your favorite HTTP library is.

Check out [the wiki](https://github.com/crazy2be/wfdr/wiki) for documentation on how the framework works.

The [HelloWorld tutorial](https://github.com/crazy2be/wfdr/wiki/HelloWorldTutorial) is a good place to go if you want to start writing a module.
