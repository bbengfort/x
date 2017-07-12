# PID

**Process ID file management for background services**

> The pid files contains the process id (a number) of a given program. For example, Apache HTTPD may write it's main process number to a pid file - which is a regular text file, nothing more than that - and later use the information there contained to stop itself. You can also use that information (just do a `cat filename.pid`) to kill the process yourself, using `echo filename.pid | xargs kill`.
>
> &mdash; [Rafael Steil](https://stackoverflow.com/questions/8296170/what-is-a-pid-file-and-what-does-it-contain)
