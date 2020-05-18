/*
 * MIT License
 *
 * Copyright (c) 2020. aberic
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package parse

import (
	"github.com/aberic/lilydb/connector"
	"strings"
)

type show struct {
	syntaxGroup []syntax
}

func newShow() *show {
	return &show{
		syntaxGroup: []syntax{
			new(conf),
			new(database),
			new(forms),
		},
	}
}

func (s *show) name() string {
	return "show"
}

func (s *show) analysis(params []string) connector.Response {
	for _, st := range s.syntaxGroup {
		if st.name() == params[0] {
			return st.analysis(params)
		}
	}
	return connector.ResultFail(syntaxErr(strings.Join(params, " ")))
}

type conf struct {
}

func (c *conf) name() string {
	return "config"
}

func (c *conf) analysis(params []string) connector.Response {
	return nil
}

type database struct {
}

func (d *database) name() string {
	return "database"
}

func (d *database) analysis(params []string) connector.Response {
	return nil
}

type forms struct {
}

func (f *forms) name() string {
	return "forms"
}

func (f *forms) analysis(params []string) connector.Response {
	return nil
}
