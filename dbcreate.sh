#!/bin/bash

mkdir db
sqlite3 ./db/forum.db < db.sql
