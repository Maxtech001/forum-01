#!/bin/bash

sqlite3 ./db/forum.db .dump > db.sql
