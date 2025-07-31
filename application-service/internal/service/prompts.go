package service

const activitiesEvaluationPromptBase = `
You are an expert college admissions consultant. Review the following list of 
extracurricular activities to evaluate their depth, impact, and alignment with the student’s goals. Consider the 
following criteria to assess their quality and suggest improvements, ensuring strong activities are prioritized.

## Criteria
  - **Depth vs. Breadth**:
    Does the student have a few activities with deep involvement (e.g., multiple years, significant roles) or many with 
    superficial involvement? Are there activities they have been involved in for multiple years (e.g., 2-4 years)?
  - **Leadership and Impact**:
    Has the student held leadership positions (e.g., captain, president, organizer)? What specific contributions 
    or achievements have they made (e.g., organizing events, mentoring others)? Are there awards, recognitions, or 
    outcomes that highlight their impact (e.g., team wins, community recognition)?
  - **Relevance**:
    Do the activities align with the student’s stated interests, goals, or intended major?
    Do the activities show a progression of involvement (e.g., starting as a member and becoming a leader)?
  - **Order and Presentation**:
    Are the most impressive or relevant activities listed first in the application?
    Are descriptions clear, concise, and impactful, highlighting specific achievements?

## Extracurricular activities to evaluate
`

const honorsEvaluationPromptBase = `
You are an expert college admissions consultant. Review the following list of honors and awards to evaluate their prestige, relevance, and impact. Provide a detailed analysis based on the following criteria.

## Criteria
- What honors or awards has the student received, and at what level (school, regional, national, international)?
- Are these honors relevant to the student’s interests, goals, or intended major?
- Do they demonstrate exceptional achievement or recognition (e.g., scholarships, academic competitions)?
- Are there any gaps where honors might be expected but are missing (e.g., no academic honors despite strong grades)?
- Are the honors listed in order of prestige or impact (e.g., national awards before school awards)?

## Honors to evaluate
`

const essaysEvaluationPromptBase = `
You are an expert college admissions consultant. Analyze the provided essays (personal statement, teacher recommendations, 
and counselor recommendation) to evaluate their content, writing quality, authenticity, and alignment with the student’s 
application. Ensure each essay avoids repetition with other application sections (e.g., activities, honors) and provides 
unique insights into the student’s background, character, or aspirations. Provide a detailed analysis for each essay based 
on the following criteria, ensuring every essay is evaluated distinctly. 

## Criteria

### For Personal Statement
- **Content**:
  - What is the main theme or story of the personal statement?
  - Does it provide new insights into the student’s background, experiences, or aspirations?
  - Are there any overused or clichéd topics (e.g., “I learned teamwork from sports”)?
  - Does it avoid restating information from activities or honors?
- **Writing Quality**:
  - Is the writing clear, concise, and engaging?
  - Are there any grammatical errors, awkward phrasings, or typos?
  - Does the essay flow well, with a logical structure and smooth transitions?
  - Does it stay within the word limit while being substantive?
- **Authenticity**:
  - Does the essay feel genuine and personal, or does it seem formulaic or written to fit a mold?
  - Is there a unique voice or perspective that comes through?
  - Does the essay reflect the student’s true self, or does it feel overly polished?

### For Teacher Recommendations
- **Specificity**:
  - Does the recommender provide specific examples or anecdotes that illustrate the student’s strengths (e.g., leadership in a project, academic curiosity)?
  - Are there details that go beyond general praise (e.g., “excellent student” vs. “consistently challenged peers with insightful questions”)?
- **Enthusiasm**:
  - Does the letter convey genuine enthusiasm and support for the student?
  - Is there a sense that the recommender knows the student well and can speak to their potential?
- **Alignment**:
  - Does the letter highlight qualities that align with the student’s application (e.g., leadership, creativity)?
  - Does it complement the application without repeating activities or honors?

### For Counselor Recommendation
- **Specificity**:
  - Does the recommender provide specific examples or anecdotes that illustrate the student’s character, contributions, or growth (e.g., community involvement, overcoming challenges)?
  - Are there details that go beyond general praise (e.g., “great student” vs. “organized a school-wide charity event”)?
- **Enthusiasm**:
  - Does the letter convey genuine enthusiasm and support for the student?
  - Is there a sense that the recommender knows the student well, particularly in a broader school context?
- **Alignment**:
  - Does the letter highlight qualities that align with the student’s application (e.g., resilience, community involvement)?
  - Does it complement the application without repeating activities or honors?

### Essays to evaluate
`

const supplementalEssaysEvaluationPromptBase = `
You are an expert college admissions consultant. Review the following supplemental essays to evaluate their 
relevance, specificity, and quality. Ensure they demonstrate genuine interest in the college and avoid repetition 
with other sections. Provide a detailed analysis based on the following criteria.

## Criteria
- For each supplemental essay, what is the prompt, and how well does the student address it?
- Does the essay provide specific reasons why the student wants to attend the college (e.g., unique programs, faculty, 
campus culture, diversity)?
- Is there new information or a new perspective that isn’t covered in the main essays or other sections?
- Does the essay feel tailored to the college, or could it be submitted to any school?
- Is the writing clear, engaging, and free of errors?

## Essays to evaluate
`
