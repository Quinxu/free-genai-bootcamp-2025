# Frontend technical specs

## Pages

### Dashboard '/dashboard'

#### Purpose
The purpose of this page is to provide a summary of learning
and act as a default page when a user visits the web-app.

#### Components
- Last Study Session
    - shows last activity used
    - shows when last activitiy used
    - summarizes wrong vs correct from last activity
    - has a link to the group
- Study Progress
    - total words study eg. 100/200
        - across all study session show the total words studied out of all possible words in our database
    - display a mastery progress eg. 50%
    - shows total words studied
    - shows total words available
- Quick Stats
    - success rate eg.80%
    - total study sessions eg. 4
    - total active groups eg. 2
    - study streak eg. 3 days
- Start Studying Button
    - goes to study activities page

#### Needed API Endpoints
- GET /api/dashboard/last_study_session
- GET /api/dashboard/study_progress
- GET /api/dashboard/quick_stats

### Study Activities '/study_activities'

#### Purpose
The purpose of this page is to show a collection of study activities with a thumbnail and its name, 
to either launch the activity or view the study activity.

#### Components
- Study Activity Card
    - shows a thumbnail of the study activity
    - the study activity name
    - a lauch button to take us to the launch page
    - the view page to view more information about past study
    sessions for this study activityhas a link to the group

#### Needed API Endpoints
- GET /api/study_activities

### Study Activity Launch '/study_activities/:id'

#### Purpose
The purpose of this page is to show the details of a study activity and its past study sessions.

#### Components
- Name of study activity
- Thumbnail of study activity
- Description of study activity
- Launch button
- Study Activities Paginated List
    - id
    - activity name
    - group name
    - start time
    - end time (inferred by the last word_review_item submitted)
    - number of review items
 
#### Needed API Endpoints
- GET /api/study_activities/:id
- GET /api/study_activities/:id/study_sessions

### Study Activity Launch '/study_activities/:id/launch'

#### Purpose
The purpose of this page is to launch a study activity.

#### Components
- Name of study activity
- Launch form
    - select field for group
    - launch now button

## Behavior
After the form is submitted, a new tab opens with the study activity based on its URL provided in the database. 

Also after form is submitted the page will redirect to the study session show page.

#### Needed API Endpoints
- POST /api/study_activities

### Words '/words'

#### Purpose
The purpose of this page is to show all words in our database.

#### Components
- Paginated Word List
    - Columns
        - Chinese
        - English
        - Correct Count
        - Wrong Count
    - Pagination with 100 items per page
    - Clicking the Chinese word will take us to the word show page.

#### Needed API Endpoints
- GET /api/words

### Word show '/words/:id'

#### Purpose
The purpose of this page is to show information about a specific word.

#### Components
- Chinese
- English
- Study Statistics
    - Correct Count
    - Wrong Count
- Word Groups 
    - show an a series of pills eg. tags
    - when group name is clicked it will take us to the group show page

#### Needed API Endpoints
- GET /api/words/:id

### Word Groups '/groups'

#### Purpose
The purpose of this page is to show a list groups in our database.

#### Components
- Paginated Group List
    - Columns
        - Group Name
        - Word Count
    - Clicking the group name will take us to the group show page.

#### Needed API Endpoints
- GET /api/groups