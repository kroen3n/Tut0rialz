
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

  I expect of you to take a look at it before diving into this tutorial... </i>
  
<i>Don't skip: </i>

<i> https://man7.org/linux/man-pages/man8/setcap.8.html   </i>
  
<i> https://man7.org/linux/man-pages/man8/getcap.8.html   </i>

<i> https://man7.org/linux/man-pages/man1/capsh.1.html    </i>




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
I will be using following Golang program - this program <a href="https://raw.githubusercontent.com/kroen3n/Tut0rialz/master/LinuX/Capabilities/write_into_file.go"> write_into_file.go </a> will add and append a couple of lines:

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

Let's try to rename the file ... We could use "mv" command, but should we now? Let's practice more Golang, with program <a href="https://raw.githubusercontent.com/kroen3n/Tut0rialz/master/LinuX/Capabilities/rename_me.go">rename_me.go</a>:

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

Since this is an ownership related challenge, my documentation mentions following:

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
Copy /bin/chown into hue's home folder, /home/hue <i> 
(pay attention from here on! this is just as an example for capabilities!! 
!!! Do not copy around tools/utilities that are not meant to be run by non-root users!!!) </i>

```
hue@kroen3n:~$ cp /bin/chown .
hue@kroen3n:~$ ls -ltr /home/hue/chown
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

As expected, no output ...

Let's apply the ownage change we need with setcap:

```
hue@kroen3n:~$  /sbin/setcap cap_chown+ep chown 
unable to set CAP_SETFCAP effective capability: Operation not permitted
hue@kroen3n:~$ 
```

Ah, 'securiteh' ... 

As a work around, ust use sudo or (if you had issues with /etc/sudoers file) just go back to root user,  
and apply the previous /sbin/setcap line.

Once applied, when you run /sbin/getcap, you should see the following output:

```
hue@kroen3n:~$ pwd
/home/hue
hue@kroen3n:~$
hue@kroen3n:~$ /sbin/getcap ./chown
chown = cap_chown+ep
```
<i> I ran 'pwd' command just to remind you that this is the chown utility we copied under /home/hue folder. 
We applied that CAP_CHOWN change only to /home/hue/chown utility </i>

And let's apply this newly changed chown on our file:

```
hue@kroen3n:~$ 
hue@kroen3n:~$ ./chown hue:hue hielau.txt
hue@kroen3n:~$
hue@kroen3n:~$ ls -ltr hielau.txt
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

And obviously, for the Linux commands fans, append a new line with “echo”:

```
hue@kroen3n:~$ more hielau.txt 
hiya
hiya again
hue@kroen3n:~$ echo "potato is Vodka" >> hielau.txt 
hue@kroen3n:~$  more hielau.txt 
hiya
hiya again
potato is Vodka
hue@kroen3n:~$  
```

<i><b> Got more? </b></i>


You want more, eh?

Take a look at ...
<br>
... following code:   https://man7.org/tlpi/code/online/dist/cap/cap_text.c.html  [1]

... and /proc filesystem documentation:  https://man7.org/linux/man-pages/man5/proc.5.html [2] 

You might want to check the <i> capability.h </i> file, from your linux host:


```
root@kroen3n:/home/hue# find / -name capability.h
/usr/include/linux/capability.h
root@kroen3n:/home/hue# 
```

From /proc filesystem documentation, following section will interest you (for now):

```
 /proc/[pid]/status 
 [...]
        * CapInh, CapPrm, CapEff: Masks (expressed in hexadecimal) of
                capabilities enabled in inheritable, permitted, and effec‐
                tive sets (see capabilities(7)).
		
		[ ... some bounding and ambient caps following ... ]
		
 [...]
````
The inheritable, permitted and effective sets are explained at [1]

<i> Time to practice </i>

Let's grab a random Linux process (in my case 1824), 

```
root@kroen3n:/home/hue# ps -ef
UID        PID  PPID  C STIME TTY          TIME CMD
[...]
root      1824     0  0 16:54 pts/1    00:00:00 bash
[...]
```


and see its capabilities:

```
root@kroen3n:/home/hue# more /proc/1824/status | grep -i cap
CapInh:	00000000a80c25fb
CapPrm:	00000000a80c25fb
CapEff:	00000000a80c25fb
CapBnd:	00000000a80c25fb
CapAmb:	0000000000000000
```

Remember, there is a tool we haven't used it yet, capsh (although, the documentation link was provided) 

<br>

Let's locate it:
```
root@kroen3n:/home/hue# whereis capsh
capsh: /sbin/capsh
```
... and use it to decode the capabilities:

```
root@kroen3n:/home/hue# /sbin/capsh --decode=00000000a80c25fb
0x00000000a80c25fb=cap_chown,cap_dac_override,cap_fowner,cap_fsetid,cap_kill,cap_setgid,cap_setuid,cap_setpcap,cap_net_bind_service,
cap_net_raw,cap_sys_chroot,cap_sys_ptrace,cap_mknod,cap_audit_write,cap_setfcap

```
Beautiful, eh?

<br></br>

Now, let's practice with a bit of Golang code. 

Suppose you want to test a small golang TCP Listener, <a href="https://github.com/kroen3n/Tut0rialz/blob/master/LinuX/Capabilities/tcp_ln.go">tcp_ln.go</a>

```
root@kroen3n:/home/hue# more tcp_ln.go
package main

import (
    "fmt"
    "net"
    "os"
)


func main() {

    // create server
    // Documentation: https://golang.org/pkg/net/
    portService := ":80"
    listener, err := net.Listen("tcp", portService)

    if err != nil {
                fmt.Fprintln(os.Stdout, err)
                os.Exit(2)
    }
    

    fmt.Println("Listening...")

    for {
        // Listening for incoming connection.
	//
        conn, err := listener.Accept()

        if err != nil {
	    fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
            os.Exit(1)
        }
	
        // Handle connections 
        go handleConnection(conn)
	
    }
}

// function to handle incoming requests 
func handleConnection(handleconn net.Conn) {

  defer handleconn.Close()

  handleconn.Write([]byte("'twas a success!Bye!\n"))

  handleconn.Close()

}

```
Just to test it, as root, build it, and run it:

```
root@kroen3n:/home/hue# go build -o tcp_ln tcp_ln.go
root@kroen3n:/home/hue#
root@kroen3n:/home/hue# ./tcp_ln
Listening...
```
Open another terminal, and see if something runs on port 80:

```
root@kroen3n:/home/hue# lsof -i :80
COMMAND   PID USER   FD   TYPE  DEVICE SIZE/OFF NODE NAME
tcp_ln   27373 root    3u  IPv6 3394583      0t0  TCP *:80 (LISTEN)
root@kroen3n:/home/hue#
```
From same new terminal, let's create a connection with our Listener on port 80  (this runs locally, of course!)
```
 root@kroen3n:/home/hue# nc localhost 80
'twas a success!Bye!
```

Let's become again "hue" user, and try to run the program:

```
 root@kroen3n:/home/hue#  su - hue
 hue@kroen3n:/home/hue$ ./tcp_ln
 listen tcp :80 : bind: permission denied 
```
Let's change the ownership, and run it again:

```
hue@kroen3n:/home/hue$  sudo chown hue:hue ./tcp_ln
hue@kroen3n:/home/hue$  ls -ltr tcp_ln
-rwxr-xr-x 1 hue hue 2945376 aug 25 16:38 tcp_ln
hue@kroen3n:/home/hue$  ./tcp_ln
listen tcp :80: bind: permission denied
```
Same error. Time to recheck the documentation for Linux capabilities. 
You will notice there are a few for networking:

```
 CAP_NET_ADMIN
              Perform various network-related operations:
              * interface configuration;
              * administration of IP firewall, masquerading, and accounting;
              * modify routing tables;
              * bind to any address for transparent proxying;
              * set type-of-service (TOS);
              * clear driver statistics;
              * set promiscuous mode;
              * enabling multicasting;
              * use setsockopt(2) to set the following socket options:
                SO_DEBUG, SO_MARK, SO_PRIORITY (for a priority outside the
                range 0 to 6), SO_RCVBUFFORCE, and SO_SNDBUFFORCE.

       CAP_NET_BIND_SERVICE
              Bind a socket to Internet domain privileged ports (port
              numbers less than 1024).

       CAP_NET_BROADCAST
              (Unused)  Make socket broadcasts, and listen to multicasts.

       CAP_NET_RAW
              * Use RAW and PACKET sockets;
              * bind to any address for transparent proxying.
```

Since we are dealing with a "binding" error, and our port is less than 1024, we should try the CAP_NET_BIND_SERVICE 

Apply capability:
```
 hue@kroen3n:/home/hue$ sudo /sbin/setcap cap_net_bind_service=+ep tcp_ln
 hue@kroen3n:/home/hue$
 hue@kroen3n:/home/hue$ /sbin/getcap tcp_ln
 tcp_ln = cap_net_bind_service+ep

 ```
 And run again as non-root user:
 
 ```
 hue@kroen3n:/home/hue$ ./tcp_ln
Listening...
```
Nice! It works! 

[... in progress...]



