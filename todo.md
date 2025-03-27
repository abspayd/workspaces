## To Do

### Logging
I want to support logging errors, especially since this is a program that typically
runs without a user-facing shell

### Structure

- root
    - load all project directories from workspace layout config json to a global

- add
    - Add a new directory that stores projects (e.g. ~/Developer)
    - Add --stow-dir flag
        - treat all subdirectories as GNU Stow projects, and add the symlink sources
        to the project global
        - E.g.:
            - ~/.dotfiles/nvim -> ~/.dotfiles/nvim/.config/nvim -- follows GNU symlink path to source
            - ~/.dotfiles/zsh -> ~/.dotfiles/.zsh/              -- there is no depth to this project, so just link to the source
        
