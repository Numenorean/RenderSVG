# Usage:
```
Make request to the http://localhost:26497/svg2png?width=450&height=150&r=0&g=0&b=0&a=255&addBackground=true
Where width and height are render size settings(must multiply or devide by the same value)
r,g,b - colour channels must be from 0 to 255(read wiki)
a - alpha channel, don't change it if you don't know what it does
addBackground - if false then result image is without background
background defines with r,g,b,a parameters
Main svg data is transmitted via base64 encoding in json format e.g.
{"svg":"long base64 string..."}
Script returns also json e.g.
{"data":"long base64 string..."}
where data is converted svg data to the png format, also in base64
```
### Credits to:
https://stackoverflow.com/a/63227777/13266817
