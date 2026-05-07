alter table entries drop constraint if exists entries_mood_chk;

update entries set mood = 'bright' where mood in ('happy', 'energized');
update entries set mood = 'tender' where mood = 'love';
update entries set mood = 'heavy' where mood = 'sad';

alter table entries
  add constraint entries_mood_chk
  check (mood in ('bright','tender','calm','heavy','stormy'));
