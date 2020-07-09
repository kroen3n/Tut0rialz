
# Linux capabilities & D0cker containers

```
 For the purpose of performing permission checks, traditional UNIX
       implementations distinguish two categories of processes: privileged
       processes (whose effective user ID is 0, referred to as superuser or
       root), and unprivileged processes (whose effective UID is nonzero).
       Privileged processes bypass all kernel permission checks, while
       unprivileged processes are subject to full permission checking based
       on the process's credentials (usually: effective UID, effective GID,
       and supplementary group list).

       Starting with kernel 2.2, Linux divides the privileges traditionally
       associated with superuser into distinct units, known as capabilities,
       which can be independently enabled and disabled.  Capabilities are a
       per-thread attribute.
```

<i> Source: https://man7.org/linux/man-pages/man7/capabilities.7.html ... 
  I expect you take a look at it before diving into this tutorial... </i>
  
<i>Also: </i>

<i> https://man7.org/linux/man-pages/man8/setcap.8.html </i>
  
<i> https://man7.org/linux/man-pages/man8/getcap.8.html </i>

<i> https://man7.org/linux/man-pages/man1/capsh.1.html </i>




<b> Linux host example </b>


This is a small example on why and how to use the power of capabilities

1) Create user "hue" and group "hue"

```
root@kroen3n:/home/hue# useradd hue
root@kroen3n:/home/hue# mkdir -p /home/hue
root@kroen3n:/home/hue# chown -R hue:hue /home/hue
```

Checking...
```
root@kroen3n:/home/hue# cat /etc/passwd | grep hue
hue:x:1000:1000::/home/hue:/bin/sh
root@kroen3n:/home/hue# 
root@kroen3n:/home/hue# su - hue
$ bash
hue@kroen3n:~$ pwd
/home/hue
```

Under user hue's home directory, as root, create an empty file:

```
root@kroen3n:/home/hue# touch hiya.txt
root@kroen3n:/home/hue#
root@kroen3n:/home/hue# ls -ltr hiya*
-rw-r--r-- 1 root root    0 Jul  9 13:18 hiya.txt
```

Now, become hue user, and let's start playing with that file:

```
root@kroen3n:/home/hue# su - hue
$ bash
hue@kroen3n:~$ ls -ltr hiya*
-rw-r--r-- 1 root root    0 Jul  9 13:18 hiya.txt
```

Suppose I want to write into that file. 
I will be using following Golang program - this program will add and append a couple of lines:

```
package main

import (
  "os"
	"io/ioutil"
	"log"
)

func main(){
	err := ioutil.WriteFile(os.Args[1], []byte("hiya\n"), 0644)

	if err != nil{
		log.Fatal(err)
	}

	file, err := os.OpenFile(os.Args[1], os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil{
		log.Println(err)
	}

	defer file.Close()

	if _, err := file.WriteString("hiya again\n"); err != nil{
		log.Fatal(err)
	}

}
```

Run the program with the name of the file as argument

```
hue@kroen3n:~$ go run write_into_file.go hiya.txt 
2020/07/09 13:24:15 open hiya.txt: permission denied
exit status 1
```

Hmm... 

Let's try to rename the file ... We could use "mv" command, but should we now? Let's practice more Golang:

```
package main

import (
	"log"
	"os"
)

func main(){
	actualFile:= os.Args[1]
	newFile := os.Args[2]

	err := os.Rename(actualFile, newFile)

	if err != nil {
		log.Fatal(err)
	}
}

```
Run the program with the name of the file and the chosen name,  as arguments:

```
hue@kroen3n:~$  go run rename_me.go hiya.txt hiya.py
hue@kroen3n:~$  
hue@kroen3n:~$  ls -ltr *.py
-rw-r--r-- 1 root root 0 Jul  9 13:31 hiya.py
hue@kroen3n:~$  
hue@kroen3n:~$ 
```

This worked. 

Back to my (now renamed) file:

```
hue@kroen3n:~$  ls -ltr *.py
-rw-r--r-- 1 root root 0 Jul  9 13:31 hiya.py
```

As you can well notice, I am trying to write into a file that is owned by root user/root group ... while I am merely a non-root user:

```
hue@kroen3n:~$  id
uid=1000(hue) gid=1000(hue) groups=1000(hue)
```

What options do I have? since this is a ownage related challenge, my documentation mentions following:

```
       CAP_CHOWN
              Make arbitrary changes to file UIDs and GIDs (see chown(2)).
```











