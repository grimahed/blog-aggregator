You'll need Postgres and Go installed to run this
You'll also need goose to run migrations.
To get started, in this order:
  - (LINUX/WSL ONLY) sudo apt update
  - refer to this to download go https://go.dev/doc/install
  - Install goose to run migrations, you will *need* to up migrate in order to use the database.
    - go install github.com/pressly/goose/v3/cmd/goose@latest
  - Install Postgres:
    -  MAC: install Postgres with: brew install postgresql@15 (you'll need v15 or later for this)
    - LINUX/WSL: sudo apt install postgresql postgresql-contrib
  - psql --version to make sure it worked
  - on LINUX/WSL, run: sudo passwd postgres to update the postgres password, write it down somewhere.
  - start the server:
      - MAC: brew services start postgresql@15
      - Linux/WSL: sudo service postgresql start
  - Use the psql shell for simplicity, it's the default client for Postgres.
      - MAC: psql postgres
      - LINUX/WSL: sudo -u postgres psql
  - In the terminal, you should see "postgres=#"
      - Create a database with CREATE DATABASE gator; (Yes, with the semicolon)
  - Connect to it with \c gator
  - Set the user password with ALTER USER postgres PASSWORD 'your-password-here';
      - You can set it to whatever you want, but you can use postgres for simplicity's sake, I did.
  - SELECT version(); to see what version of Postgres you're running.
  - You can exit by typing either exit or \q

Install this with: go install github.com/grimahed/blog-aggregator@latest.

To set up your config file so you can run the gator command in the prompt (needs to know the string to make a connection):
  - Make a file in your home directory (type cd ~ then press enter if you aren't sure what "home directory" means)
    (Type nano .gatorconfig.json and press enter to enter edit mode of a blank file. touch "your-file-name-here" creates an empty file,
    nano both creates and sends you into it, if it already exists, it just puts you into edit mode)
  - paste this inside and change "your-set-password" to the one you have set for the user:
    {"db_url":"postgres://postgres:your-set-password@localhost:5432/gator?sslmode=disable",
     "current_user_name":"YOU"
    }
  - up migrate with goose
    - MAC: goose "postgres://your-username:@localhost:5432/gator" up
    - LINUX/WSL: goose "postgres://postgres:password-you-set@localhost:5432/gator" up
  - The migration string before "up" is your connection string
    - you can connect to the database with psql <connection-string>
    - if you decide to down migrate from 1 to 0, it will wipe the entire database to nothing, *BE CAREFUL*
      - simple as as up migrating again to use it, but I wouldn't recommend it.
      - if you just want to reset the database to a clean slate, it's better to run the reset command (See below) 

A few commands you can use (everything is case-sensitive)

register - adds a user to the users table
login - logs in the user, needs the name of said user.
feed - adds a feed to the feeds table, needs a name (configurable), and the url.
follow - follows a feed  (will obviously want to be a user who *didn't* make the feed)
unfollow - It, well. Lets the user unfollow a feed. Needs the name of the feed.
agg - collects feeds and creates posts of said feeds. This takes a timer argument. (i.e 30s, 1m, 1hr, etc.) Try not to DOS them. 
      This is an intentionally designed infinite loop, so you'll need CTRL+C to cancel out the command and keep track of how many you've collected/how
      long it's been running.
browse - allows you to browse posts that have been collected by agg, the amount is configurable. (so, say, you want 4 posts. put 4 in as your argument).
reset - resets the entire database to a fresh clean slate
