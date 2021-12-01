<div id="top"></div>

[![Contributors][contributors-shield]][contributors-url] [![Forks][forks-shield]][forks-url] [![Stargazers][stars-shield]][stars-url] [![Issues][issues-shield]][issues-url] [![MIT License][license-shield]][license-url] [![LinkedIn][linkedin-shield]][linkedin-url]

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/S4ND1X/PacmanGo">
    <img src="images/logo.png" alt="Logo" width="140" height="140">
  </a>

  <h3 align="center">Golang Pacman</h3>

  <p align="center">
    Eat the dots and avoid the ghosts!
    <br />
    <a href="https://github.com/S4ND1X/PacmanGo"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/S4ND1X/PacmanGo">View Demo</a>
    ·
    <a href="https://github.com/S4ND1X/PacmanGo/issues">Report Bug</a>
    ·
    <a href="https://github.com/S4ND1X/PacmanGo/issues">Request Feature</a>
  </p>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#architecture">Architecture</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->

## About The Project

![Golang Pacman Screen Shot][product-screenshot]

<br/>

![Winning Screen Shot][winning-screenshot]

Golang is a programming language that is used to build software and is used to develop web applications, one of the advantages of using Golang is that it is a statically typed language and it is easy to write code.

Also handling multi threading is a great feature of Golang and it is very easy to use this wy we decided to create a Pacman game using Golang.

Also pacman is a famous game in the world and it is very easy to play this game but the work of creating a Pacman game is very interesting and challenging since you get to touch a lot of aspects from the language such inpute handling, graphics, multi threading, etc.

<p align="right">(<a href="#top">back to top</a>)</p>

### Built With

This project was pretty simple in terms of technology. The only dependency was the Go language

- [Golang](https://golang.org/)
- [Ubuntu](https://www.ubuntu.com/)
- [Git](https://git-scm.com/)

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- GETTING STARTED -->

## Getting Started

It's easy to get started with this project, simply follow the steps below and you'll be up and running in no time, the hard part is winning the game!

### Prerequisites

Please make sure you have the following installed before you start

- **Golang**

Ubuntu:

```sh
  sudo apt-get install golang-go
```

Mac:

```sh
  brew install golang
```

Windows:

```sh
  https://golang.org/dl/
```

<br/>

- **Visual Studio Code**

Ubuntu:

```sh
  sudo apt-get install code
```

Mac:

```sh
    brew cask install visual-studio-code
```

Windows:

```sh
  https://code.visualstudio.com/download
```

### Installation

_Don't worry if this is your first time using Golang this, guide should set everything you need_

1. Clone the repository and navigate to root directory

```sh
  git clone git@github.com:S4ND1X/PacmanGo.git
  cd PacmanGo
```

2. Install Golang dependencies

```sh
   make build
```

3. Run the application

```sh
   go run main.go
```

4. To specify number of enemies specify a integer value in the command line between 1 and 12

```sh
   go run main.go < 1 - 12 >
```

### Build and Run

- Build and run the application

```sh
   make run enemies=< 1 -12 >
```

- Test the application

```sh
   make test
```

- Build the application

```sh
   make build
```

- Install dependencies

```sh
   make deps
```

- Clean previous build

```sh
   make clean
```

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- USAGE EXAMPLES -->

## Usage

This is a simple Pacman game, you can use the arrow keys to move the Pacman and the WASD keys to move the ghosts.

If you want to play with more enemies, you set the number of enemies in the command line.

    go run main.go < 1 - 12 >

If you die terminal process will be ended and you can start again by running the application again.

    go run main.go

If you want to end the game press the `esc` key.

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- Architecture -->

## Architecture

We try to keep our architecture simple and focused by dividing our project into logical sections and keeping each section focused on a single task.

### Logical Sections

<div align="center">

![Architecture Screenshot][architecture-screenshot]

</div>

### Flow of Events

<div align="center">

![Flow Diagram Screenshot][flow_diagram-screenshot]

</div>

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- LICENSE -->

## License

Distributed under the MIT License. See `LICENSE.txt` for more information.

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- CONTACT -->

## Contact

- Jorge Sanchez Diaz - @S4ND1X - georgesanchez.diazjr@gmail.com
- Jose Luis Aguilar Nucamendi - @JosxLuis - luisaguilarnucamendi@gmail.com
- Agustin Salvador Quintanar De la Mora - @AgusQuintanar - agusquintanar17@gmail.com

Project Link: [https://github.com/S4ND1X/PacmanGo](https://github.com/S4ND1X/PacmanGo)

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- ACKNOWLEDGMENTS -->

## Acknowledgments

We would like include the following resource that helped us create this project:

- [Golang](https://golang.org/): The programming language that powers the internet.
- [Pacgo](https://github.com/danicat/pacgo): A Golang guide to creating a Pacman game.
- [Readme Template](https://github.com/othneildrew/Best-README-Template): An awesome README template to jumpstart your projects!

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->

[contributors-shield]: https://img.shields.io/github/contributors/S4ND1X/PacmanGo.svg?style=for-the-badge
[contributors-url]: https://github.com/S4ND1X/PacmanGo/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/S4ND1X/PacmanGo.svg?style=for-the-badge
[forks-url]: https://github.com/S4ND1X/PacmanGo/network/members
[stars-shield]: https://img.shields.io/github/stars/S4ND1X/PacmanGo.svg?style=for-the-badge
[stars-url]: https://github.com/S4ND1X/PacmanGo/stargazers
[issues-shield]: https://img.shields.io/github/issues/S4ND1X/PacmanGo.svg?style=for-the-badge
[issues-url]: https://github.com/S4ND1X/PacmanGo/issues
[license-shield]: https://img.shields.io/github/license/S4ND1X/PacmanGo.svg?style=for-the-badge
[license-url]: https://github.com/S4ND1X/PacmanGo/blob/master/LICENSE.txt
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://www.linkedin.com/in/jorgesanchezdiaz/
[product-screenshot]: images/screenshot.jpeg
[architecture-screenshot]: images/architecture.png
[flow_diagram-screenshot]: images/flow_diagram.png
[winning-screenshot]: images/winning.jpeg
