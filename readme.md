# Present

Is a TUI for presenting slides...

## Doc

### The .pres file

```
config as yaml
~~~

Slide 1

---

Slide 2

---

@comment{
    Wow should be bold ....
}
Slide 3 is @b{wow}!


```


### Config

The presentation can/should begin with the following (it is yaml syntax)

```yaml
# custom styles can optionally be set
style:
    # name of the style (used like other named styles)
    # the not all options (bold, italic, bg, fg) must be set
    fancy:
        # defines the foreground color
        fg: '#04B575'
        # defines the background color
        bg: '#00000'
        # ...
        bold: true
        # ..
        italic: true

# the author can optionally be defined
author: Sojamann

# This is the most important as this seperates the
# config from the slides
~~~

```

### Blocks

- *@code[type]{...}*                        syntax highlighting of ...
- *@note{...}*                              full width text block
- *@warning{...}*                           full width text block which stands out
- *@comment{...}*                           comment which is not rendered
- *@img[{width: .., height: ..}]{ /path }*  renderes image

### Named Styles

- *!h{..}*      heading
- *!b{..}*      bold
- *!i{..}*      italic
- *!red{..}*    red
- *!green{..}*  green
- *!yellow{..}* yellow
- *!blue{..}*   blue
- *!white{..}*  white

## ToDo
[ ] themes
[ ] layouts (split left and right)?
[ ] tables
[ ] plugins?
[ ] code execution?

