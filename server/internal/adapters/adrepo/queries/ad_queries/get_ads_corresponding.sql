select id, title, text, author_id, published, created, last_modified
from ads
where published in $1 and
      created in $2 and
      author_id in $3 and
      title in $4