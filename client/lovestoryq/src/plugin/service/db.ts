function db(momentList){ 
    const db = openDatabase('hoper', '1.0', 'hoper DB', 2 * 1024 * 1024)
    db.transaction(function (tx) {
      tx.executeSql('CREATE TABLE IF NOT EXISTS Moments (id unique, content)')
    })
   
    db.transaction(function (tx) {
      for (const item of momentList) {
        tx.executeSql(
          `INSERT INTO Moments (id, content) VALUES (${item.id},'${
            item.content
          }')`
        )
      }
    })

    db.transaction((tx)=> {
      tx.executeSql(
        'SELECT * FROM Moments',[],
        (tx, results) => {
          const len = results.rows.length
          console.log(`查询记录条数:${len}`)
          for (let i = 0; i < len; i++) {
            console.log(results.rows.item(i).content)
          }
        },
        null
      )
    })
}