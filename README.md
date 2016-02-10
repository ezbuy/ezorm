# ezorm

ezorm is an code-generation based ORM lib for golang, supporting mongodb/mysql/sql server.

data model is defined with YAML file:

```yaml
Product:
  db: mongo
  fields:
    - name: CarTeamOpinion
      type: int
    - name: DeputyDirectorOpinion
      type: m_voiture.ApprovalSuggestion
    - name: DirectorOpinion
      type: m_voiture.ApprovalSuggestion
```
