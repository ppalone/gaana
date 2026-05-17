package gaana

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_decAESCBCPKCS(t *testing.T) {
	message := "5sUy0ScJH6MdO8l6y8C5FhdX77t4XJMk+FkpgtwIWA+tRheLiAHgH7T9RHqb20pLF7UtFSp3QQ2YQ/KpjHyl8+pMA6m9uHKJ92H53yYL/NoU1bXT7eiLChSquOdMdCNx2vy6P9J8H3/23Fi+8ZLiJs6SuQlgxfvmmX9oeyeRpQg7gsije5Yh6UEAkkZtPdtO+PRY4YmmoQhA9Fa+bwN2+8gg502d9Ut9znevBsyDFb2Faq8MRGdcNdfEE61WOHUu3hZx4bnriPxJNVlTZ21ZAqEnAXICeDbUT0XyIcggJifgRf/X+eZS6VEXKBQDwmfs="
	decrypted, err := decAESCBCPKCS(message)

	assert.NoError(t, err)
	assert.NotEmpty(t, decrypted)
	assert.True(t, strings.HasPrefix(decrypted, "https://"))
	assert.Contains(t, decrypted, ".m3u8")
}
