# GORM Find Example

This is an example project that demonstrates how to use GORM's Find method to 
find a record by one of it's associations fields.


## Overview

In this project we have Users and Orders.  Users have a name and Orders have a 
description and belong to a User.


## Use Case

I'm using GORM in an API project where I need to filter my Orders by various 
fields, including fields that exist in related tables.


# Running the Example

This example project uses Docker.  To run the example simply run the following
command:

```
docker-compose run --rm example
```

When the project runs it will first drop and recreate the users and orders 
tables, create some test records, and then run a few queries.
