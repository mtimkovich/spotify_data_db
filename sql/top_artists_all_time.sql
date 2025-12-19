select artist, sum(ms_played) / 60000
from songs
group by artist
order by sum(ms_played) desc
limit 40;
