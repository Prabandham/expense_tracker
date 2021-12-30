{{define "groupedCredits"}}
SELECT sum(q1.amount) AS "amount", name AS "type" FROM
(SELECT c.amount, ct.name FROM credits AS c
LEFT JOIN credit_types AS ct
ON cASt(c.credit_type_id AS uuid) = ct.id WHERE c.user_id = ?) AS q1
GROUP BY q1.name
{{end}}

{{define "groupedDebits"}}
SELECT sum(q1.amount) AS "amount", name AS "type" FROM 
(SELECT c.amount, ct.name FROM debits AS c
LEFT JOIN debit_types AS ct
ON cASt(c.debit_type_id AS uuid) = ct.id WHERE c.user_id = ?) AS q1
GROUP BY q1.name
{{end}}