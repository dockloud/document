package document

import "time"

type Document struct {
	ID        int
	OwnerID   string
	Name      string
	Content   string
	IsPrivate bool
	IsDeleted bool
	CreatedAt time.Time
	EditedAt  time.Time
}

func (d *Document) UpdateName(name string) error {
	err := ValidateName(name)
	if err != nil {
		return err
	}
	d.Name = name
	return nil
}

func (d *Document) UpdateContent(content string) error {
	err := ValidateContent(content)
	if err != nil {
		return err
	}
	d.Content = content
	return nil
}

func (d *Document) SwitchPrivacy(privacy bool) {
	if privacy == d.IsPrivate {
		return
	}
	d.IsPrivate = privacy
}
