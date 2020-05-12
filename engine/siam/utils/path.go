/*
 * MIT License
 *
 * Copyright (c) 2020 aberic
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

package utils

import (
	"github.com/aberic/gnomon"
	"github.com/aberic/lilydb/config"
	"path/filepath"
)

// PathFormIndexFile 表索引文件路径
//
// databaseID 数据库唯一id
//
// formID 表唯一id
//
// indexID 表索引唯一id
func PathFormIndexFile(databaseID, formID, indexID string) string {
	return gnomon.StringBuild(config.Obtain().DataDir, string(filepath.Separator), databaseID, string(filepath.Separator), formID, string(filepath.Separator), indexID, ".idx")
}

// PathFormFile 表文件路径
//
// databaseID 数据库唯一id
//
// formID 表唯一id
func PathFormFile(databaseID, formID string) string {
	return filepath.Join(config.Obtain().DataDir, databaseID, formID, "form.dat")
}
