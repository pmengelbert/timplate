## Timplate
A simple tool that converts a timesheet written in yaml to a nice-looking PDF timesheet, written in golang.  Uses the `pdflatex` tool to compile a PDF.

### Dependencies
Install the `pdflatex` tool, which is necessary for compiling into PDF

### Basic installation
```
git clone https://github.com/pmengelbert/timplate.git
cd timplate
go build .
sudo cp timplate /usr/local/bin
```

### Basic usage
You'll need a yaml file in the following format:
```
name: Peter Engelbert
rate: 32
startDate: Feb 03, 2020
endDate: Feb 13, 2020
records:
    - date: 02/03/20
      times: 
          - 1030-1915
      hours: 8.75
      description: 
        - In office
    - date: 02/04/20
      times: 
        - 1030-1830
      hours: 8
      description: 
            - deployed postgres server
            - created new endpoints to add users to the database, remove users, and show all users
            - made changes for OCI PR
```

Then, simply run:
```
timplate <filename.yaml>
```
This will generate `timesheet.pdf`
