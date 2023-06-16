# Halocat - BE

This github repository is the backend application for Halocat App.

## About Halocat

The "Hallocat" application is designed to facilitate cat lovers in consulting with veterinarians. With online consultation features, we aim to ensure that their cats' health is well-maintained and they have easy access to the necessary medical care."

## API Documentation

You can find complete API documentation at the following link: https://documenter.getpostman.com/view/18409946/2s93m1bR3M

## How To Run On Local Machine

Requirements:
- Go v1.19.4 or later
    - How to install on Windows: https://www.youtube.com/watch?v=kxD8p-aPYzM
    - How to install on Linux (Ubuntu): https://www.youtube.com/watch?v=mNMaBXFY_4Y
    - How to install on Mac: https://www.youtube.com/watch?v=fPjcp48dpPM

If you really want to run it yourself on your local machine, you can follow these steps:

1. ```git
   git clone -b master https://github.com/rickyromansyah2045/halocat-backend-go.git
    ```
2. ```
   cd halocat-backend-go
    ```
3. ```
   go mod download
    ```
4. ```
   cp .env.example .env
    ```

5. ```
   go run main.go -production=false
    ```
6. Enjoy and Happy Explore! 

## Tech Stack

- Go Using Gin Framework

## Contacts

- https://www.linkedin.com/in/ricky-romansyah-47831518b
