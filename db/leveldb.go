//
// Copyright 2022 YumeMichi. All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.

// This binary provides sample code for using the gopacket TCP assembler and TCP
// stream reader.  It reads packets off the wire and reconstructs HTTP requests
// it sees, logging them.
//
package db

import (
	"errors"
	"sifrank/config"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type LevelDbImpl struct {
	ldb *leveldb.DB
}

var (
	LevelDb LevelDbImpl
	err     error
)

func init() {
	LevelDb.InitDb()
}

func (ldb *LevelDbImpl) InitDb() {
	ldb.ldb, err = leveldb.OpenFile(config.Conf.LevelDbPath, nil)
	if err != nil {
		panic(err.Error())
	}
}

func (ldb *LevelDbImpl) Get(key []byte) (res []byte, err error) {
	if len(key) == 0 {
		return nil, errors.New("key is null")
	}
	res, err = ldb.ldb.Get(key, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ldb *LevelDbImpl) Put(key, value []byte) (err error) {
	if len(key) == 0 {
		return errors.New("key is null")
	}
	err = ldb.ldb.Put(key, value, nil)
	if err != nil {
		return err
	}
	return nil
}

func (ldb *LevelDbImpl) List() (res map[string]string) {
	res = make(map[string]string)
	iter := ldb.ldb.NewIterator(nil, nil)
	for iter.Next() {
		res[string(iter.Key())] = string(iter.Value())
	}
	return res
}

func (ldb *LevelDbImpl) ListPrefix(prefix []byte) (res map[string]string) {
	res = make(map[string]string)
	iter := ldb.ldb.NewIterator(util.BytesPrefix(prefix), nil)
	for iter.Next() {
		res[string(iter.Key())] = string(iter.Value())
	}
	return res
}
