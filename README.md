# Groupie Tracker Filters

## 📚 Introduction
Groupie Trackers is a website that uses data from an API to diplay information about different bands. The information includes name(s), image, when they began their activity, the date of their first album, the members and their tour dates and locations. The dates and locations are diplayed on a map.

Groupie Tracker Filters consists on letting the user filter the artists/bands that will be shown.

Project must incorporate at least these four filters:
- filter by creation date
- filter by first album date
- filter by number of members
- filter by locations of concerts

Filters must be of at least these two types:
- a range filter (filters the results between two values)
- a check box filter (filters the results by one or multiple selection)

## 👟 Requirements to run

- Go >= 1.18
- Only standard Go packages were used
- Bash terminal window

## 🏃‍♂️ Running the program

Just use 
`go run .`
Then open the link http://localhost:8080/

## 🧪 Testing the program
Audit can be found [here](https://github.com/01-edu/public/blob/master/subjects/groupie-tracker/filters/audit.md)

## ✏️ Notes
The server is written in Go. HTML, CSS and JavaScript are used for frontend. JavaScript is used only for displaying tour infomation on the map.

## 🤴 Authors
@Brooklyn_95

@kretesaak

@margus.aid