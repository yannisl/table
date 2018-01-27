# The Caption Package

The Go caption package offers routines and datastructures for setting and manipulating the
captions of figures, tables and similar.

A caption provides a short explanation, or description accompanying an illustration, photograph or table. The structure and semantics of a caption are provided by a label and text. 

```
  Figure 12:  This is a figure.
```

In the above example the Figure is the label, and the text is "This is a figure."

```html
<figure>
  <img src="img_pulpit.jpg" alt="The Pulpit Rock" width="304" height="228">
  <figcaption>Fig1. - A view of the pulpit rock in Norway.</figcaption>
</figure>
```

## Label

The label has a descriptive name and possibly a number and can have a terminating symbol.

## Presentation

I have tried as far as possible to split the presentation from the semantic form. The initial
focus is on LaTeX2e.  