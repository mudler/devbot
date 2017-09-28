package bot

import "github.com/boltdb/bolt"

func DBRemoveSingleKey(bucket string, key string) bool {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open(DB, 0600, nil)
	if err != nil {
		return false
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {

		bkt, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		return bkt.Delete([]byte(key))
	})

	if err == nil {
		return true
	}
	return false
}

func DBListKeys(bucket string) ([]string, error) {
	var list []string
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open(DB, 0600, nil)
	if err != nil {
		return list, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {

		boltbkt, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		boltbkt.ForEach(func(k, v []byte) error {
			list = append(list, string(k))
			return nil
		})
		return nil
	})

	return list, err
}

func DBListValues(bucket string) ([]string, error) {
	var list []string
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open(DB, 0600, nil)
	if err != nil {
		return list, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {

		boltbkt, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		boltbkt.ForEach(func(k, v []byte) error {
			list = append(list, string(v))
			return nil
		})
		return nil
	})

	return list, err
}

func DBAllKeyValue(bucket string) (map[string]string, error) {
	var list map[string]string
	list = make(map[string]string)
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open(DB, 0600, nil)
	if err != nil {
		return map[string]string{}, err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {

		boltbkt, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		boltbkt.ForEach(func(k, v []byte) error {
			list[string(k)] = string(v)
			return nil
		})
		return nil
	})

	return list, err
}

func DBGetSingleKey(bucket string, ket string) (string, error) {
	var answer []byte
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open(DB, 0600, nil)
	if err != nil {
		return "", err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {

		bucketbolt, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}

		// Retrieve if already exists
		if v := bucketbolt.Get([]byte(ket)); v != nil {
			answer = v
		}

		return err
	})

	return string(answer), err
}

func DBPutKey(bucket string, key string) bool {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open(DB, 0600, nil)
	if err != nil {
		return false
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {

		bkt, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}

		return bkt.Put([]byte(key), nil)
	})
	if err == nil {
		return true
	}
	return false
}

func DBPutKeyValue(bucket string, key string, value string) bool {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open(DB, 0600, nil)
	if err != nil {
		return false
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {

		bkt, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}

		return bkt.Put([]byte(key), []byte(value))
	})
	if err == nil {
		return true
	}
	return false
}
