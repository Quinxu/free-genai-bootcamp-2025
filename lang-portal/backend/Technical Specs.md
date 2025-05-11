# Backend Server Technical Specs

## Business Goal 
Andrew's application will do the following:
- A language learning school wants to build a prototype of learning portal which will act as three things:
- Inventory of possible vocabulary that can be learned, 
- Act as a  Learning record store (LRS), providing correct and wrong score on practice vocabulary
- A unified launchpad to launch different learning apps.
 
My application will do the following:
- learn a language by providing an english setence, then ask the user to provide the corresponding sentence in the target language.
- With help function, the application provides the sentence structure in the target language, and the likely used words.
- After an user enteres the sentence, the application will provide the feedback on the correct and wrong words. 

## Technical Requirements
- The backend will be built using Go
- The database will be SQLite3
- The API will be built using Gin
- The API will always return JSON
- There will be no authentication or authorization
- evething will be treated as a single user

## Database schema

Our database will be a single sqlite database called `words.db` that will be in the root of the project folder of `backend_go`

We have the following tables:
- words - stored vocabulary words
    - id (primary key) integer
    - chinese string
    - english string   
    - parts json
    - created_at datetime    
- words_groups - join table between words and groups many to many
    - id (primary key) integer
    - word_id (foreign key) integer
    - group_id (foreign key) integer
    - created_at datetime
- groups - thematic groups of words
    - id (primary key) integer
    - name string
    - created_at datetime
- study_sessions - records of study sessions grouping word_review_items
    - id (primary key) integer
    - word_id (foreign key) integer
    - group_id (foreign key) integer
    - study_activity_id (foreign key) integer
    - study_session_id (foreign key) integer
    - correct boolean
    - created_at datetime
- study_activities - a specific study activity, linking a study session to group
    - id (primary key) integer
    - study_session_id (foreign key) integer
    - group_id (foreign key) integer
    - created_at datetime
- word_review_items - records of word practice, determining if the word was correct or not
    - id (primary key) integer
    - word_id (foreign key) integer
    - study_session_id (foreign key) integer
    - correct boolean
    - created_at datetime


## API Endpoints

### GET /api/dashboard/last_study_session
Returns information about the most recent study session.

#### JSON Response
```json
{
  "id": 123,
  "group_id": 456,
  "created_at": "2025-02-08T17:20:23-05:00",
  "study_activity_id": 789,
  "group_id": 456,
  "group_name": "Basic Greetings"
  }
}
```

### GET /api/dashboard/study_progress
Returns study progress statistics.
Please note that the frontend will determine progress bar based on total words studied and total available words.

#### JSON Response
```json
{
  "total_words_studied": 3,
  "total_available_words": 124
}
```

### GET /api/dashboard/quick_stats
Returns quick statistics about study progress.

#### JSON Response
```json
{
  "success_rate": 80.0,
  "total_study_sessions": 4,
  "total_active_groups": 3,
  "study_streak_days": 4
}
```

### GET /api/study_activities/:id
Returns a single study activity by ID.

#### JSON Response
```json
{
  "id": 1,
  "name": "Vocabulary Quiz",
  "thumbnail_url": "https://example.com/thumbnail.jpg",
  "description": "Practice your vocabulary with flashcards"
}
```
// Possible error:
```json
{
  "error": "Activity not found"
}
```

### GET /api/study_activities/:id/study_sessions
Returns paginated study sessions for a given activity.

#### JSON Response
```json
{
  "items": [
    {
      "id": 123,
      "activity_name": "Vocabulary Quiz",
      "group_name": "Basic Greetings",
      "start_time": "2025-02-08T17:20:23-05:00",
      "end_time": "2025-02-08T17:30:23-05:00",
      "review_items_count": 20
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_items": 100,
    "items_per_page": 20
  }
}
```

### POST /api/study_activities

#### required params: 
    - group_id integer 
    - study_activity_id integer

#### JSON Response
```json
{
  "id": 124,
  "group_id": 123
}
```

### Get /api/words
    - pagination with 100 items per page

#### JSON Response
```json
{
  "items": [
    {
      "chinese": "你好",
      "english": "Hello",
      "correct_count": 5,
      "wrong_count": 2
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_items": 500,
    "items_per_page": 100
  }
}
``` 

### Get /api/words/:id

#### JSON Response
```json
{
  "english": "Hello",
  "stats": {
    "correct_count": 5,
    "wrong_count": 2
  },
  "groups": [
    {
      "id": 1,
      "name": "Basic Greetings"
    }
  ]
}
```

### Get /api/groups
    - pagination with 100 items per page

#### JSON Response
```json
{
  "items": [
    {
      "id": 1,
      "name": "Basic Greetings"
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_items": 100,
    "items_per_page": 20
  }
}
```

### Get /api/groups/:id

#### JSON Response
```json
{
  "id": 1,
  "name": "Basic Greetings"
}
```

### Get /api/groups/:id/words
    - pagination with 100 items per page

#### JSON Response
```json
{
  "items": [
    {
      "id": 1,
      "chinese": "你好",
      "english": "Hello",
      "parts": {
        "part1": "part1",
        "part2": "part2"
      }
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_items": 100,
    "items_per_page": 20
  }
}
```

### Get /api/groups/:id/study_sessions

#### JSON Response
```json
{
  "items": [
    {
      "id": 123,
      "group_id": 456,
      "created_at": "2025-02-08T17:20:23-05:00",
      "study_activity_id": 789,
      "group_id": 456,
      "group_name": "Basic Greetings"
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_items": 100,
    "items_per_page": 20
  }
}
```

### GET /api/study_sessions
    - pagination with 100 items per page

#### JSON Response
```json
{
  "items": [
    {
      "id": 123,
      "group_id": 456,
      "created_at": "2025-02-08T17:20:23-05:00",
      "study_activity_id": 789,
      "group_id": 456,
      "group_name": "Basic Greetings"
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_items": 100,
    "items_per_page": 20
  }
}
```

### GET /api/study_sessions/:id

#### JSON Response
```json
{
  "id": 123,
  "group_id": 456,
  "created_at": "2025-02-08T17:20:23-05:00",
  "study_activity_id": 789,
  "group_id": 456,
  "group_name": "Basic Greetings"
}
```

### GET /api/study_sessions/:id/words

#### JSON Response
```json
{
  "items": [
    {
      "id": 1,
      "chinese": "你好",
      "english": "Hello",
      "parts": {
        "part1": "part1",
        "part2": "part2"
      }
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_items": 100,
    "items_per_page": 20
  }
}
```

### POST /api/reset_history

#### JSON Response
```json
{
  "success": true,
  "message": "Study history has been reset"
}
```

### POST /api/full_reset

#### JSON Response
```json
{
  "success": true,
  "message": "System has been fully reset"
}
```

### POST /api/study_sessions/:id/words/:word_id/review
#### Required Params 
    - id (study_session_id) integer
    - word_id integer
    - correct boolean

#### Request Payload
```json
{
  "correct": true
}
```

#### JSON Response
```json
{
  "success": true,
  "word_id": 456,
  "study_session_id": 789,
  "correct": true
  "created_at": "2025-02-08T17:20:23-05:00"
}
```

## Mage Tasks

Mage is a task runner for Go.
Let's list out possible tasks we need for our lang portal.

### Initialize Database

This task will initialize the sqlite database called `words.db`

### Migrate Database
This task will run a series of migrations sql files on the database

Migration live in the 'migrations' folder.
The migration files will be run in order of their file name.
The file names should look like this:

```sql
0001_init.sql
0002_create_words_table.sql
```


### Seed Data
This task will import json files and transform them into target data for our database.

All seed files live in the `seeds` folder.
In our task we should have DSL to specific each seed file and its expected group word name.

```json
[
  {
    "chinese": "你好",
    "english": "Hello",
    "parts": {
      "part1": "part1",
      "part2": "part2"
    }
  },
  ...
]
```


