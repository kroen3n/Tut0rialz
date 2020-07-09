
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
hue@kroen3n:~$  go run rename_me.go hiya.txt hielau.txt
hue@kroen3n:~$  
hue@kroen3n:~$  ls -ltr hielau.*
-rw-r--r-- 1 root root 0 Jul  9 13:31 hielau.txt
hue@kroen3n:~$  
hue@kroen3n:~$ 
```

This worked. 

Back to my (now renamed) file:

```
hue@kroen3n:~$  ls -ltr hielau.*
-rw-r--r-- 1 root root 0 Jul  9 13:31 hielau.txt
```

As you can well notice, I am trying to write into a file that is owned by root user/root group ... while I am merely a non-root user:

```
hue@kroen3n:~$  id
uid=1000(hue) gid=1000(hue) groups=1000(hue)
```

What options do I have - from a "capability" point of view? 

Since this is an ownage related challenge, my documentation mentions following:

```
       CAP_CHOWN
              Make arbitrary changes to file UIDs and GIDs (see chown(2)).
```



<i><b> Can't I just apply chown? </b></i>

Of course you can ... but as long as it's not sudo-ed!

```
hue@kroen3n:~$  chown hue:hue hielau.txt
chown: changing ownership of 'hielau.txt': Operation not permitted
hue@kroen3n:~$  
``` 

So, what's next? 

Locate chown:

```
hue@kroen3n:~$ whereis chown
chown: /bin/chown
hue@kroen3n:~$
```
Copy it into hue's home folder <i> (pay attention from here on! this is just as an example for capabilities!! 
Do not copy around tools/utilities that are not meant to be run by non-root users!!!) </i>

```
hue@kroen3n:~$ cp /bin/chown .
hue@kroen3n:~$ ls -ltr chown
-rwxr-xr-x 1 hue  hue  72512 Jul  9 14:20 chown
```
Let's check where are my CAP tools - <i> getcap</i> - that will provide the capabilities that are already set-up, 
and <i> setcap </i> - that will (obviously) set-up the capabilities I require.
```
hue@kroen3n:~$ whereis getcap
getcap: /sbin/getcap
hue@kroen3n:~$ 
hue@kroen3n:~$ whereis setcap
setcap: /sbin/setcap
hue@kroen3n:~$ 
```
Let's see how to use them:

```
hue@kroen3n:~$  
hue@kroen3n:~$  /sbin/getcap /home/hue/chown
hue@kroen3n:~$ 
```

As expected, nothing to see there ...

Let's apply the ownage change we need with setcap:

```
hue@kroen3n:~$  /sbin/setcap cap_chown+ep chown 
unable to set CAP_SETFCAP effective capability: Operation not permitted
hue@kroen3n:~$ 
```

Ah, 'securiteh' ... 

Just use sudo or (if you had issues with /etc/sudoers file) just go back to root user,  and apply the previous /sbin/setcap line

Once it's done, when you run /sbin/getcap, you should see the following output:

```
hue@kroen3n:~$ pwd
/home/hue
hue@kroen3n:~$
hue@kroen3n:~$ /sbin/getcap ./chown
chown = cap_chown+ep
```
<i> I ran 'pwd' command just to remind you that this is the chown utility we copied under /home/hue folder. 
We applied that CAP change only to /home/hue/chown </i>

And let's apply this newly changed chown on our file:

```
hue@kroen3n:~$ 
hue@kroen3n:~$ ./chown hue:hue hielau.txt
hue@kroen3n:~$
hue@kroen3n:~$ ls -ltr hiya.py 
-rw-r--r-- 1 hue hue 0 Jul  9 14:41 hielau.txt
```
You can see that now the file is owned by non-root user, hue.

Let's try to write something into our file:

```
hue@kroen3n:~$ go run write_into_file.go hielau.txt 
hue@kroen3n:~$ more hielau.txt 
hiya
hiya again
hue@kroen3n:~$
```

Yey! No more permission errors!

And obviously, for the Linux commands fans, let’s append a new line with “echo”:

```
hue@dd1c0ba95ae3:~$ more hielau.txt 
hiya
hiya again
hue@dd1c0ba95ae3:~$ echo "potato is Vodka" >> hielau.txt 
hue@dd1c0ba95ae3:~$ more hielau.txt 
hiya
hiya again
potato is Vodka
hue@dd1c0ba95ae3:~$ 
```

<i><b> Got more? </b></i>


























