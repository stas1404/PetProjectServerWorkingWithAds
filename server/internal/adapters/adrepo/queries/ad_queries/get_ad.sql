select id, title, text, author_id, published, created, last_modified
from ads
where id = $1