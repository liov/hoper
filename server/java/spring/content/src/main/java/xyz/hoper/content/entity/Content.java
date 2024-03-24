package xyz.hoper.content.entity;


import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serializable;


/**
 * @Description TODO
 * @Date 2022/11/21 15:18
 * @Created by lbyi
 */
@Table(name="content",schema="public")
@Entity
@Data
@NoArgsConstructor
@AllArgsConstructor
public class Content implements Serializable {
    @Id
    Long id;

    String name;
}