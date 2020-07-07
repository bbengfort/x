# Clock

I've had a simple `clock.py` program in my `~/bin` directory since I started programming. This little CLI utility prints the current time in the local or UTC timezone and formats it for a variety of use cases. I generally combine this script with `pbcopy` to quickly copy and paste the time into different documents.

I use this tool so much, that I thought it would be nice to extend it to be able to do simple time computations (e.g. a very common task I have is to determine the date 6 weeks from now). The issue is that my Python script has a third party dependency, namely python-dateutil for timezone support. Why not rewrite this simple helper in Go? Thus the version 2.0 clock command was born here.