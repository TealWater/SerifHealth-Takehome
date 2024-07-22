## Packages
- bufio
- bytes
- compress/gzip
- encoding/json
- fmt
- log
- os
- strings

## Go version: 1.21.5

## Design Decisions:
I decided to open the gzip file and load each byte or character into the buffer one by one. This is done by using the `buf io ` package. I decided to use `bufio.ScanBytes` func to load bytes in one by one. I could have used `bufio.ScanWords or bufio.ScanLines ` but an issue kept rearing its head. It was skipping chunks of data that were relevant to locating the URLs, so it seemed easier if I just moved arcoss the file one byte or character at a time.

The first step is to create a buffer with a size of 16 bytes, this is meant to look for the byte sequence `in_network_files`. Once we can locate this value in the file, we can then dump the data in the array into a struct and from there I can inspect the URLs.

The real question is... how do we know if the perspective link is related to NY? This was discoveralble by first taking several EINs and looking at the Anthem website: [https://www.anthem.com/machine-readable-file/search/]. If you search a few EINs you'll notice that the New York (NY) value is encoded as a number, but this number is always changing depending on the insurer. Whether it is Anthem or Blue Cross, the NY number will be different. However, there are a set of characters that remain constant and unique, depending on the insurer, that is what will be used to determine if the perspective link is associated with NY. So if you split the link by `underscore(_)` and look at the second index, you can then deteermine if the character sequence is associated with NY.

As we are gathering the links, I've noticed that the very last link in `in_network_files` is a gzip file with all of the links asscoated with its respective employer/client. There is no need to capture this and it crashes my search operation, it is discarded as the program panics and I recover it... this will be seen when `2024/07/22 02:15:49 hmm, found a link that is does not pertain to in_network_files` prints.

Once we have saved all of the valid links, all we have to do is dump them into a struct and then Marshall them into a `.json` file. There is also the issue of some characters being escaped therefore not making the links clickable if you were to redirect the output to stdout. This was remedied by creating a func to encode the values without escaping them.


## Run/Execution Time
 - 22 min

 ## Development Time:
 - 3 Days

 ## How to run program
 ```
    1)  You need to download the table of contents gzip file from the Anthem website [https://www.anthem.com/machine-readable-file/search/] and place it in the root of the project. 
        **Make sure the name matches what the program has or this will break it.**
    2) type [go run main.go] in your terminal from the root directory of this project
    3) look for the file [results.json]
 ``` 