# blog-generator

## Installation

```bash
go get github.com/eleztian/blog-generator
```

## Usage & Customization

## Configuration

Example Config File:

```yml
generator:
    repo: 'https://github.com/eleztian/blog.git'
    tmp: 'tmp'
    dest: 'www'
    npg: 15
blog:
    url: 'https://www.eleztian.xyz'
    language: 'en-us'
    description: ' -- Crazy Snail --<br/>Never stop'
    dateformat: '02.Jan.2006'
    title: 'Tab.Blog'
    author: 'Tab Eleztian'
    github: 'https://github.com/eleztian'
    frontpageposts: 10
    statics:
        files:
            - src: 'static/favicon.ico'
              dest: 'favicon.ico'
            - src: 'static/robots.txt'
              dest: 'robots.txt'
            - src: 'static/css/vec.css'
              dest: 'css/vec.css'
            - src: 'static/css/vec.css.map'
              dest: 'css/vec.css.map'
            - src: 'static/js/highlight.min.js'
              dest: 'js/highlight.min.js'
            - src: 'static/welcome.jpg'
              dest: 'welcome.jpg'
        templates:
            - src: 'static/welcome.html'
              dest: ''
```