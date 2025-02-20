---
id: gqsr69tq4uqy97i9v02vfk0
title: Gptproject
desc: ''
updated: 1740028987533
created: 1740028313298
---

## Certified Argo Project Associate dump generator

### Project files

- [CAPA Curriculum](https://github.com/cncf/curriculum/blob/master/CAPA_Curriculum.pdf)
- for each concatenated doc of category since there is a limit for the file number
  - `find docs -type f -name "*.md" | grep -v "[A-Z]" | grep -v "README\|CONTRIBUTION" | xargs cat > path/to/kubecadet/capa/{category}-concat-docs.md`
  - [argo cd docs](https://github.com/flavono123/argo-cd/tree/master/docs)
  - [argo workflows docs](https://github.com/flavono123/argo-workflows/tree/master/docs)
  - [argo rollouts docs](https://github.com/flavono123/argo-rollouts/tree/master/docs)
  - [argo events docs](https://github.com/flavono123/argo-events/tree/master/docs)

### Instructions(Prompt)

[!NOTE]

```plaintext
generate questions for mock capa exam
when i ask the category in the exam
reference the file uploaded  {category}-concat-docs.md and CAPA_Curriculum.pdf

first, generate multiple-choice 10 questions about the category what i asked without solutions
and then i solved i will let you know my solution, check that

when submit my solutions, answer with revision and some reading material from concat doc or official docs published on the internet for incorrect problems.

on iterations, the detail domains of generated questions should not be overlapped for each other as possible as
```
