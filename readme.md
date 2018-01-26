# Table

Table is a Go package that translates .csv files to nice LaTeX2e tables. Currently it is
very idiosyncratic and the API will change as it evolves.

Go is well suited for the development of Command Line Interfaces as well as the manipulation of text files. It's built-in library provides off-the-shelf libraries for parsing encoded files such as .csv or .json files.

LaTeX has packages that can handle csv data files directly, but for larger files they are limited and tend to slow compilation.

With the package one can export from excel to csv and then use a Go preprocessor to build up the tables. The tables are saved to disk and can then be imported to LaTeX with the `input{<tablename.tex>}` command.

