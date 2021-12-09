# Welcome to Bookshelf - Backend 📖

![bookshelf-logo-github](https://user-images.githubusercontent.com/76471929/145391946-8870d37b-fab8-4fd4-8a68-000d33d02d15.png)

 ## About Bookshelf 📚
 
 Bookshelf is a smart-booking application for efficiently saving and using bookmarks while working in the browser. To get started head on over to [Bookshelf web client](https://web-bookshelf.netlify.app) and Sign up for an account.

## How to add a custom search engine 📑
This example uses Google Chrome, however setup is similar across most browsers.

1. First, find the settings page.
2. Next, locate the Search engine tab on the left hand side.
3. Click Manage search engines.
4. In the Other search engines section, click Add.
 - Undersearch engine, choose a name; e.g. Bookshelf.
 - Under Keyword, choose a keyword to invoke Bookshelf; e.g. bk, shelf, etc.
 - Under URL, copy and paste your unique URL.

## Get started developing 🖥️
This is the repository for the backend. If you would like to work on the frontend, check out the [frontend repository](https://github.com/conalli/bookshelf-client-web) 📘.

The backend is written entirely in Go, using Redis and MongoDB (with MongoDB Atlas). 
To get started
 - Clone this repository.
 - To build the project and run the server, in the terminal run:
 ``` 
 docker-compose build
 docker-compose up -d
 ```
  - When you are done, run:
 ```
 docker-compose down
 ```
