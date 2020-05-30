/*Employee 表包含所有员工信息，每个员工有其对应的工号 Id，姓名 Name，工资 Salary 和部门编号 DepartmentId 。

+----+-------+--------+--------------+
| Id | Name  | Salary | DepartmentId |
+----+-------+--------+--------------+
| 1  | Joe   | 85000  | 1            |
| 2  | Henry | 80000  | 2            |
| 3  | Sam   | 60000  | 2            |
| 4  | Max   | 90000  | 1            |
| 5  | Janet | 69000  | 1            |
| 6  | Randy | 85000  | 1            |
| 7  | Will  | 70000  | 1            |
+----+-------+--------+--------------+
Department 表包含公司所有部门的信息。

+----+----------+
| Id | Name     |
+----+----------+
| 1  | IT       |
| 2  | Sales    |
+----+----------+
编写一个 SQL 查询，找出每个部门获得前三高工资的所有员工。例如，根据上述给定的表，查询结果应返回：

+------------+----------+--------+
| Department | Employee | Salary |
+------------+----------+--------+
| IT         | Max      | 90000  |
| IT         | Randy    | 85000  |
| IT         | Joe      | 85000  |
| IT         | Will     | 70000  |
| Sales      | Henry    | 80000  |
| Sales      | Sam      | 60000  |
+------------+----------+--------+
解释：

IT 部门中，Max 获得了最高的工资，Randy 和 Joe 都拿到了第二高的工资，Will 的工资排第三。销售部门（Sales）只有两名员工，Henry 的工资最高，Sam 的工资排第二。
*/

SELECT
    d.Name AS 'Department', e1.Name AS 'Employee', e1.Salary
FROM
    Employee e1
        JOIN
    Department d ON e1.DepartmentId = d.Id
WHERE
        3 > (SELECT
                 COUNT(DISTINCT e2.Salary)
             FROM
                 Employee e2
             WHERE
                     e2.Salary > e1.Salary
               AND e1.DepartmentId = e2.DepartmentId
    );

/*若是在8.0以上版本，可直接调用内置函数实现排名，且效率要更高一些：

row_number()，连续排名，同分不同名，例如90、80、80、70，排名后是1、2、3、4
rank()，同分同名，但是跳级，例如90、80、80、70，排名后是1、2、2、4
dense_rank()，致密排名，同分同名，不跳级，例如90、80、80、70，排名后是1、2、2、3
本题即为dense_rank()的情况。
*/
SELECT
    c.name as Department, c.Employee, c.salary as Salary
FROM
    (
        SELECT
            b.name, a.name as Employee, a.salary,
            dense_rank()OVER(PARTITION BY a.departmentid ORDER BY a.salary DESC) AS rn
        FROM employee a LEFT JOIN department b ON a.departmentid=b.id
        ORDER BY 1,4,2
    )c
WHERE c.rn IN (1,2,3)
  AND c.name IS NOT NULL
ORDER BY 1,3 DESC,2 DESC;

/*将当前表按部门和salary从高到低进行排序，则同一部门的所有薪资集中，并且是按从高到低排序
自定义三个变量：@pD，表示上一个部门id，初始化为null；@pS，表示上一个薪资，初始化为null，当然这里的上一个薪资可能是同部门也可能是不同部门；@r，表示当前排名，初始化为0（虽然MySQL不区分大小写，但我们在书写中仍然加以区分方便识读。另外，这里选择用一个虚拟表的方式进行变量初始化）
执行以下逻辑：
A. 如果当前部门id与前一部门id相同（@pD=DepartmentId），则排名信息是在当前部门的基础上进行更新，又区分两种情况：
与前一薪资相同（@pS=salary），则排名不变（@r:=@r）
否则排名+1（@r:=@r+1）
B. 如果当前部门id与前一部门id不同，则说明排名需重新开始，所以排名赋值为1（@r:=1）
*/
SELECT
    d.NAME department, t.NAME employee, salary
FROM
    ( SELECT
          *, @r := IF(@pD = departmentid, IF(@pS = salary, @r, @r + 1 ), 1 ) AS 'rank',
          @pD := departmentid,
          @pS := salary
      FROM
          employee, ( SELECT @pS := NULL, @pD := NULL, @r := 0 ) init
      ORDER BY
          departmentid, salary DESC ) t
        JOIN department d ON t.departmentid = d.id
WHERE
        t.rank <=3;