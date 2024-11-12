# bootdev_gator
boot.dev learning project. building a rss aggregator.

## Requirements
This program can be run from your command line. You need Go and Postgres installed.

## Config
You need to create a `.gatorconfig.json` file in your home directory with the following structure:
````json
{
    "db_url": "postgres://username:@localhost:5432/database?sslmode=disable"
}
````
Replace the values with your database connection string.

## Usage
After you installed the program, you have to register a user, add at least one feed and start the aggregator. This will let the program run from your command line. To stop aggregating you have to force stop (ctrl+c).

## Commands
The following commands are avaible and can be typed in your command line:
- create a new user: ``gator register <username>``
- login user: ``gator login <username>``
- list all user: ``gator users``
- add a feed: ``gator addfeed <name> <url>``
- list all feeds: ``gator feeds``
- follow a feed (already in DB): ``gator follow <url>``
- unfollow a feed: ``gator unfollow <url>``
- start the aggregator: ``gator agg 1m``
- view posts: ``gator browse [limit]``


Have Fun!